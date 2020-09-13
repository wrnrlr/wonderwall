import crel from 'crelt'
class LoginForm extends HTMLElement {
    constructor() {
        super()
        this.$message = crel('div', {})
        this.$email = crel("input", {id: 'email', type: 'email'})
        this.$button = crel("button", {}, 'Reset Password')
    }
    connectedCallback() {
        this.$button.addEventListener('click', _ => console.log('click btn'))
        this.appendChild(crel('label', {'for': 'email'})); this.appendChild(this.$email)
        this.appendChild(this.$button)
    }
    async submit() {

    }
}

document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('login-form', LoginForm)
})
