import crel from 'crelt'



class Editor extends HTMLElement {
    constructor() { super() }
    connectedCallback() {}
    disconnectedCallback() {}
}

class Tools extends HTMLElement {
    get value() { return this.getAttribute('value') }
    set value(v) { this.setAttribute('value', v) }
    constructor() {
        super()
        this.tools = ['selection', 'pen', 'text', 'image', 'shape']
    }
    connectedCallback() {
        this.tools.forEach(t => this.appendChild(crel('div', {value: t, onclick: this.fireValueEvent})))
        this.value = this.tools[0]
    }
    attributeChangedCallback(name, oldValue, newValue) {
        if (name === 'value') {
            this.childNodes[this.tools.indexOf(oldValue)].classList.remove('active')
            this.childNodes[this.tools.indexOf(newValue)].classList.remove('active')
            this.value = newValue
        }
    }
    fireValueEvent(e) {
        console.log(e.target.getAttribute('value'))
        this.dispatchEvent(new CustomEvent('value', {detail: e.target.getAttribute('value')}))
    }
}

class Application extends HTMLElement {
    constructor() {
        super()
    }
    connectedCallback() {
        this.$editor = crel('wall-editor')
        this.$tools = crel('wall-tools')
        this.$tools.addEventListener('value', this.toolValue)
        this.appendChild(this.$editor)
        this.appendChild(this.$tools)
    }
    toolValue(e) {
        console.log('new tool value: ' + e.detail)
    }
}

document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('wall-tools', Tools)
    window.customElements.define('wall-editor', Editor)
    window.customElements.define('wall-app', Application)
})
