import { h, Component } from 'preact';
import axios from 'axios'; // Assuming you use axios for HTTP requests. If not, adjust accordingly.

class Login extends Component {
    state = {
        username: '',
        password: '',
        errorMessage: ''
    };

    handleInputChange = (e) => {
        this.setState({ [e.target.name]: e.target.value });
    };

    handleLogin = async () => {
        try {
            const formData = new FormData();
            formData.append('username', this.state.username);
            formData.append('password', this.state.password);

            const response = await axios.post('http://localhost:8080/login', formData);
            // Store the JWT token in localStorage
            localStorage.setItem('token', response.data.token);

            // Inform the parent (App) about successful login
            if (this.props.onSuccess) {
                this.props.onSuccess();
            }

            // Navigate to the main page
            window.location.href = '/';

        } catch (error) {
            this.setState({ errorMessage: 'Login failed. Please try again.' });
        }
    };

    render() {
        return (
            <div>
                <input
                    type="text"
                    name="username"
                    placeholder="Username"
                    value={this.state.username}
                    onChange={this.handleInputChange}
                />
                <input
                    type="password"
                    name="password"
                    placeholder="Password"
                    value={this.state.password}
                    onChange={this.handleInputChange}
                />
                <button onClick={this.handleLogin}>Login</button>
                {this.state.errorMessage && <p>{this.state.errorMessage}</p>}
            </div>
        );
    }
}

export default Login;
