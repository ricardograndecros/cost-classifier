import { h, Component } from 'preact';

export default class TransactionCard extends Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedLabel: this.props.transaction.label_id || null
        };
    }

    handleLabelChange = (e) => {
        const newLabel = e.target.value;
        this.setState({ selectedLabel: newLabel });
        this.props.onUpdate(this.props.transaction.ID, newLabel);
    }

    render({ transaction, labels }, { selectedLabel }) {
        return (
            <div className="transaction-card">
                <div className="transaction-info">
                    <h2>{transaction.Title}</h2>
                    <p>{transaction.Amount.toFixed(2)} (Paid from {transaction.AccountNumber || 'N/A'})</p>
                </div>
                <div className="transaction-actions">
                    <select value={selectedLabel} onLabelChange={this.handleLabelChange}>
                        {labels.map(label => (
                            <option value={label.ID}>{label.Name}</option>
                        ))}
                    </select>
                </div>
            </div>
        );
    }
}
