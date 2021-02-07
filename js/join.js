import crel from 'crelt'
class JoinForm extends HTMLElement {
    constructor() {
        super()
        this.$message = crel('div', {})
        this.$email = crel("input", {id: 'email', type: 'email'})
        this.$password = crel("input", {id: 'password', type: 'password'})
        this.$terms = crel("input", {id: 'terms', type: 'checkbox'})
        this.$button = crel("button", {'class': 'green'}, 'Join')
    }
    connectedCallback() {
        this.$button.addEventListener('click', _ => console.log('click btn'))
        this.appendChild(crel('label', {'for': 'email'}, 'Email')); this.appendChild(this.$email)
        this.appendChild(crel('label',  {'for': 'password'}, 'Password')); this.appendChild(this.$password)
        this.appendChild(this.$terms); this.appendChild(crel('label',  {'for': 'terms'}, 'Accept ', crel('a', {href:'/terms'}, "Terms & Conditions")))
        this.appendChild(this.$button)
    }
    async submit() {

    }
}

document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('join-form', JoinForm)
})
