import crel from 'crelt'
class LoginForm extends HTMLElement {
    constructor() {
        super()
        this.$message = crel('div', {})
        this.$email = crel("input", {id: 'email', type: 'email'})
        this.$password = crel("input", {type: 'password'})
        this.$button = crel("button", {'class': 'blue'}, 'Login')
    }
    connectedCallback() {
        this.$button.addEventListener('click', _ => console.log('click btn'))
        this.appendChild(crel('label', {'for': 'email'}, 'Email')); this.appendChild(this.$email)
        this.appendChild(crel('label',  {'for': 'password'}, 'Password')); this.appendChild(this.$password)
        this.appendChild(this.$button)
    }
    async submit() {

    }
}

document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('login-form', LoginForm)
})
