import { h, Component } from 'preact';
import TransactionCard from './TransactionCard';
import Popup from './Popup';
import LabelForm from './LabelForm';

export default class TransactionList extends Component {
    constructor(props) {
        super(props);
        this.state = {
            transactions: [],
            labels: [],
            isPopupOpen: false,
            activeTab: 'transactions', // transactions or labels
            isDeletePopupOpen: false,
            labelToDelete: null
        };
    }

    componentDidMount() {
        this.loadTransactions();
        this.loadLabels();
    }

    toggleDropdown = () => {
        this.setState(prevState => ({ isDropdownOpen: !prevState.isDropdownOpen }));
    }

    selectLabel = (labelName) => {
        this.setState({ selectedLabel: labelName });
    }



    loadTransactions = async () => {
        try {
            const token = localStorage.getItem('token');  // Retrieve the JWT token from localStorage
            const response = await fetch('http://localhost:8080/transactions', {
                headers: {
                    'Authorization': `Bearer ${token}`  // Include the JWT token in the Authorization header
                }
            });

            if (!response.ok) {
                throw new Error('Failed to load transactions');
            }

            const transactions = await response.json();
            if (!Array.isArray(transactions)) {
                throw new Error('Received data is not an array');
            }

            this.setState({ transactions });
        } catch (error) {
            console.error('Error loading transactions:', error);
        }
    };



    loadLabels = async () => {
        try {
            const token = localStorage.getItem('token');  // Retrieve the JWT token from localStorage
            const response = await fetch('http://localhost:8080/labels', {
                headers: {
                    'Authorization': `Bearer ${token}`  // Include the JWT token in the Authorization header
                }
            });

            if (!response.ok) {
                throw new Error('Failed to load labels');
            }

            const labels = await response.json();
            if (!Array.isArray(labels)) {
                throw new Error('Received data is not an array');
            }

            this.setState({ labels });
        } catch (error) {
            console.error('Error loading labels:', error);
        }
    };

    handleLabelChange = (event) => {
        const selectedLabelName = event.target.value;
        this.setState({ selectedLabel: selectedLabelName });
    }


    handleUpdateLabel = async (oldLabelName, updatedLabelData) => {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8080/labels/edit?old_name=${oldLabelName}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updatedLabelData)
        });
        if (response.ok) {
            this.loadLabels();  // Reload labels to ensure the list is updated
        } else {
            console.error('Error updating label:', await response.text());
        }
    }



    openPopup = () => {
        this.setState({ isPopupOpen: true });
    }

    closePopup = () => {
        this.setState({ isPopupOpen: false });
    }

    handleCreateLabel = async (labelData) => {
        const token = localStorage.getItem('token');
        const encodedColor = encodeURIComponent(labelData.color);  // Encode the color value
        const response = await fetch(`http://localhost:8080/labels/new?name=${labelData.name}&color=${encodedColor}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (response.ok) {
            await this.loadLabels();  // Reload labels to ensure the list is updated
            this.setState({isPopupOpen: false});
        } else {
            console.error('Error creating label:', await response.text());
        }
    }


    handleDeleteLabel = async (labelName) => {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8080/labels/delete?name=${labelName}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (response.ok) {
            await this.loadLabels();  // Reload labels to ensure the list is updated
        } else {
            console.error('Error deleting label:', await response.text());
        }
    }

    openDeleteConfirmation = (labelName) => {
        // Logic to open the Popup component and store the label to delete
        this.setState({ isDeletePopupOpen: true, labelToDelete: labelName });
    }


    handleDeleteLabelConfirmation = async() => {
        const token = localStorage.getItem('token');
        fetch(`http://localhost:8080/labels/delete?name=${encodeURIComponent(this.state.labelToDelete.Name)}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        })
            .then(response => {
                if (response.ok) {
                    // Refresh the labels list or remove the label from state
                    this.loadLabels()
                } else {
                    // Handle any errors that occur during the fetch
                    console.error('Error deleting label:', response.statusText);
                }
            })
            .catch(error => console.error('Fetch error:', error));

        this.closeDeleteConfirmation()
    }

    closeDeleteConfirmation = () => {
        this.setState({ isDeletePopupOpen: false, labelToDelete: null });
    }

    handleTabChange = (tabName) => {
        this.setState({ activeTab: tabName });
    }

    render(_, { transactions, labels, isPopupOpen, activeTab }) {
        return (
            <div className="transaction-list">
                <div className="tabs">
                    <button
                        className={`tab-button ${activeTab === 'transactions' ? 'active' : ''}`}
                        onClick={() => this.handleTabChange('transactions')}
                    >
                        Transactions
                    </button>
                    <button
                        className={`tab-button ${activeTab === 'labels' ? 'active' : ''}`}
                        onClick={() => this.handleTabChange('labels')}
                    >
                        Labels
                    </button>
                </div>

                <div className={`tab-content ${activeTab === 'transactions' ? 'active' : ''}`}>
                    {transactions.map(transaction => (
                        <TransactionCard
                            transaction={transaction}
                            labels={labels}
                            onLabelChange={this.handleLabelChange}
                        />
                    ))}
                </div>

                <div className={`tab-content ${activeTab === 'labels' ? 'active' : ''}`}>
                    {labels.map(label => (
                        <div className={`label-card`} style={{ backgroundColor: label.Color }}>
                            <span>{label.Name}</span>
                            <span className="delete-icon" onClick={() => this.openDeleteConfirmation(label)}>
                                <i className="fas fa-trash-alt"></i>
                            </span>
                        </div>
                    ))}

                    {/* "Create new label" card */}
                    <div className="label-card create-label-card" onClick={this.openPopup}>
                        <span>Create new label</span>
                        <span className="plus-symbol">+</span>
                    </div>
                </div>

                <Popup isOpen={isPopupOpen} onClose={this.closePopup}>
                    <LabelForm onCreateLabel={this.handleCreateLabel} />
                </Popup>

                <Popup isOpen={this.state.isDeletePopupOpen} onClose={this.closeDeleteConfirmation}>
                    <p>Are you sure you want to delete the label "{this.state.labelToDelete?.Name}"?</p>
                    <button onClick={this.handleDeleteLabelConfirmation}>Yes</button>
                    <button onClick={this.closeDeleteConfirmation}>No</button>
                </Popup>

            </div>
        );
    }
}
