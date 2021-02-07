import crel from 'crelt'
const emailRegex = /[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/
function validEmail(s) { return emailRegex.test(s) }
function validPassword(s) { return s.length >= 8 }
class LoginForm extends HTMLElement {
    constructor() {
        super()
        this.$message = crel('div', {class: 'message'}, "message")
        this.$email = crel("input", {id: 'email', type: 'email'})
        this.$password = crel("input", {type: 'password'})
        this.$button = crel("button", {'class': 'login'}, 'Login')
    }
    connectedCallback() {
        this.$button.addEventListener('click', _ => this.submit())
        this.appendChild(this.$message)
        this.appendChild(crel('div', {class: 'fields'},
            crel('label', {'for': 'email'}, 'Email'), this.appendChild(this.$email),
            crel('label',  {'for': 'password'}, 'Password'), this.appendChild(this.$password)))
        this.appendChild(this.$button)
    }
    async submit() {
        console.log("submit login")
        const email = this.$email.value, password = this.$password.value
        while (this.$message.lastChild) this.$message.removeChild(this.$message.firstChild)
        if (!validEmail(email)) { this.$message.appendChild(crel('div', {}, "invalid email")) }
        if (!validPassword(password)) { this.$message.appendChild(crel('div', {}, "invalid password, should be at least 8 characters long")) }
        if (this.$message.childElementCount !== 0) return
        console.log("valid login form, submitting")
        const form = {email, password}
        try {
            const resp = await fetch('/login', {method: 'POST', redirect: 'follow', body: JSON.stringify(form)})
            if (resp.status === 400) {
                console.log('login bad request')
            }
            if (resp.redirected) {
                console.log('redirect')
                // window.location.href = resp.url;
            }
        } catch (e) {
            console.log("error fetch login")
        }
    }
}

document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('login-form', LoginForm)
})
