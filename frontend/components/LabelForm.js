import { h, Component } from 'preact';

export default class LabelForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            name: '',
            color: '#000000'
        };
    }

    handleInputChange = (e) => {
        const { name, value } = e.target;
        this.setState({ [name]: value });
    }

    handleSubmit = (e) => {
        e.preventDefault();
        this.props.onCreateLabel(this.state);
    }

    render(_, { name, color }) {
        return (
            <form onSubmit={this.handleSubmit}>
                <div>
                    <label for="name">Name:</label>
                    <input
                        type="text"
                        id="name"
                        name="name"
                        value={name}
                        onChange={this.handleInputChange}
                    />
                </div>
                <div>
                    <label for="color">Color:</label>
                    <input
                        type="color"
                        id="color"
                        name="color"
                        value={color}
                        onChange={this.handleInputChange}
                    />
                </div>
                <button type="submit">Create Label</button>
            </form>
        );
    }
}
