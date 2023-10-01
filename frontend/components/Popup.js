import { h, Component } from 'preact';

export default class Popup extends Component {
    render({ isOpen, onClose, children }) {
        if (!isOpen) {
            return null;
        }

        return (
            <div className="popup">
                <div className="popup-inner">
                    <button onClick={onClose} className="popup-close-btn">
                        X
                    </button>
                    {children}
                </div>
            </div>
        );
    }
}
