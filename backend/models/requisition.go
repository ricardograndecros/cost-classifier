package models

type Requisition struct {
	RequisitionId string `gorm:"primary_key"`
	Username      string `gorm:"primary_key"`
	AgreementId   string `gorm:"not null"`
	Reference     string `gorm:"not null"`
	BankId        string `gorm:"not null"`
}
