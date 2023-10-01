package fetcher

import (
	"errors"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"cost-classifier/backend/config"
	"cost-classifier/backend/db"
	"cost-classifier/backend/models"
	"github.com/ricardograndecros/go-nordigen"
	"gorm.io/gorm"
)

const redirectPort = ":1234"

func NewNordigenClient() go_nordigen.Client {
	secretId := config.AppConfig.NordigenSecretId
	secretKey := config.AppConfig.NordigenSecretKey

	c, err := go_nordigen.NewClient(secretId, secretKey)
	if err != nil {
		log.Fatalf("error while creating go-nordigen client. %s", err)
	}
	return *c

}

func fetchUserRequisition(username string, c go_nordigen.Client) models.Requisition {
	var requisition models.Requisition
	result := db.DB.Where("username = ?", username).First(&requisition)
	if result.Error != nil {
		log.Print("Error querying for requisition")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Print("Requisition not found, creating one")
			// no requisitions found for the user, create a new one and then return
			reqId := CreateNewRequisition(username, c)
			log.Print("Created a requisition")
			// call database to ensure it was created
			result := db.DB.Where("requisition_id = ?", reqId).First(&requisition)
			if result.Error != nil {
				log.Print("Error while fetching newly created requisition")
			}
		} else {
			log.Fatalf("Database query error: %s", result.Error)
		}

	}
	log.Print("Returning requisition", requisition)
	return requisition
}

func RetrieveAccountTransactions(c go_nordigen.Client, username string) {
	requisition := fetchUserRequisition(username, c)

	nordigenRequisition, err := c.GetRequisitionsById(requisition.RequisitionId)
	if err != nil {
		log.Printf("Error while fetching user's bank requisitions")
		return
	}
	log.Print("Return from requisitionby id")
	for _, account := range nordigenRequisition.Accounts {
		log.Printf("First iteration: %s", account)
		//get account information
		accountInfo, err := c.GetAccountInfo(account)
		if err != nil {
			log.Printf("Error while fetching account information. %v", err)
		}
		// find most recent stored transaction
		mostRecentTransactionDate, err := FindMostRecentTransactionDate(account)
		if err != nil {
			log.Fatalf("Error while fetching most recent transaction date")
		}
		log.Print("Entra a buscar las account transactions")
		// retrieve transactions since the most recent saved transaction
		transactions, err := c.GetAccountTransactions(account,
			mostRecentTransactionDate.Format("2006-01-02"),
			time.Now().Format("2006-01-02"))
		if err != nil {
			log.Fatalf("Error fetching the transactions")
		}

		// save them to database
		log.Print("Va a guardarlas a la db")
		log.Print(transactions)
		TransactionListToDatabase(*transactions, *accountInfo)
	}

}

func TransactionListToDatabase(transactions go_nordigen.Transactions, account go_nordigen.Account) {
	log.Print("Entra a crear las transacciones")
	for _, transaction := range append(transactions.Booked, transactions.Pending...) {
		log.Print("Iterando por ", transaction)
		var existingTransaction models.Transaction
		log.Print("TRANSACTION ID: ", transaction.TransactionID)
		result := db.DB.Where("transaction_id = ?", transaction.TransactionID).First(&existingTransaction)
		log.Print(existingTransaction)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Print("Creando transacción en DB", transaction)
				amount, _ := strconv.ParseFloat(transaction.TransactionAmount.Amount, 64)

				newTransaction := models.Transaction{
					TransactionID: transaction.TransactionID,
					Title:         transaction.RemittanceInformationUnstructured,
					Amount:        amount,
					Currency:      transaction.TransactionAmount.Currency,
					AccountId:     account.Id,
					AccountIban:   account.Iban,
					Date:          time.Time{},
				}
				if err := db.DB.Create(&newTransaction).Error; err != nil {
					log.Printf("Error saving transaction to database: %v", err)
				}
			} else {
				log.Printf("Database query error: %v", result.Error)
			}
		} else {
			log.Printf("Transaction %v already exists in database", existingTransaction.TransactionID)
		}
	}
}

func FindMostRecentTransactionDate(account string) (*time.Time, error) {
	var transaction models.Transaction
	result := db.DB.Where("account_id = ?", account).Order("date desc").First(&transaction)
	if result.Error != nil {
		log.Print("Entra a crear fecha")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// no records for the account, return current month's first day
			now := time.Now()
			returnDate := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
			return &returnDate, nil
		} else {
			return nil, result.Error
		}
	}

	return &transaction.Date, nil
}

func CreateNewRequisition(username string, c go_nordigen.Client) (requisitionId string) {
	// step 2: choose a bank
	institutionId := config.AppConfig.BankID

	// step 3: create an end user agreement
	agreement, err := c.CreateUserAgreement(institutionId)

	if err != nil {
		log.Fatalf("error while creating the end user agreement. %s", err)
	}

	// step 4: build a link

	// find the current reference for the user

	r := go_nordigen.Requisition{
		InstitutionID:     agreement.InstitutionId,
		Redirect:          "http://localhost" + redirectPort,
		Agreement:         agreement.Id,
		Reference:         "0",
		Language:          "ES",
		AccountSelection:  false,
		RedirectImmediate: false,
	}

	log.Print(r)
	requisition, err := c.NewRequisition(r)
	if err != nil {
		log.Fatalf("error while creating the requisition. %s", err)
	}

	openBrowser(requisition.Link)

	waitForCallback()

	// save requisition to DB
	dbRequisition := models.Requisition{
		RequisitionId: requisition.Id,
		Username:      username,
		AgreementId:   requisition.Agreement,
		Reference:     requisition.Reference,
		BankId:        requisition.InstitutionID,
	}
	db.DB.Create(dbRequisition)

	return dbRequisition.RequisitionId
}

// Open URL in default browser
func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Fatalf("Failed to open browser with url %s: %v", url, err)
	}
}

func waitForCallback() {
	callbackSignal := make(chan bool)

	http.HandleFunc(redirectPort, func(w http.ResponseWriter, r *http.Request) {
		// Signal that we have received the callback
		callbackSignal <- true

		// Respond to the HTTP request
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<p>Link to account created successfully.</p>
					<p>Redirecting to your cost-classifier...</p>
					<script>
						setTimeout(function() {
							window.location.href = 'http://localhost:1234';
							window.close(); // This might not always work due to browser security policies
						}, 3000);
					</script>
				</body>
			</html>
		`))
	})

	go func() {
		http.ListenAndServe(":8888", nil)
	}()

	// Wait for the signal
	<-callbackSignal
}
