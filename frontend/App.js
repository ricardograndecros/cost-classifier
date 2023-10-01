import { h, Component } from 'preact';
import TransactionList from './components/TransactionList';
import Login from './components/Login';  // Import the Login component

export default class App extends Component {
    state = {
        isLoggedIn: !!localStorage.getItem('token')  // Check if token exists as a simple way to verify if user is logged in
    };

    handleLoginSuccess = () => {
        this.setState({ isLoggedIn: true });
    };

    render() {
        return (
            <div id="app">
                {this.state.isLoggedIn ? <TransactionList /> : <Login onSuccess={this.handleLoginSuccess} />}
            </div>
        );
    }
}
