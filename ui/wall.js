import crel from 'crelt'
class Editor extends HTMLElement {
    constructor() { super() }
    connectedCallback() {}
    disconnectedCallback() {}
}
class Tools extends HTMLElement {
    static get observedAttributes() { return ['value'] }
    get value() { return this.getAttribute('value') }
    set value(v) { this.setAttribute('value', v) }
    constructor() { super(); this.tools = ['selection', 'pen', 'text', 'image', 'shape'] }
    connectedCallback() {
        this.tools.forEach(t => this.appendChild(crel('div', {value: t, onclick: this.fireValueEvent})))
        this.value = this.tools[0]
    }
    attributeChangedCallback(name, oldValue, newValue) {
        if (name === 'value') {
            if (oldValue) this.childNodes[this.tools.indexOf(oldValue)].classList.remove('active')
            this.childNodes[this.tools.indexOf(newValue)].classList.add('active')
        }
    }
    fireValueEvent(e) { this.dispatchEvent(new CustomEvent('value', {detail: e.target.getAttribute('value'), bubbles: true})) }
}
class Application extends HTMLElement {
    connectedCallback() {
        this.$editor = crel('wall-editor')
        this.$tools = crel('wall-tools', {})
        this.$tools.addEventListener('value', e => this.toolValue(e))
        this.appendChild(this.$editor)
        this.appendChild(crel('div', {class:'options menu'}))
        this.appendChild(this.$tools)
    }
    toolValue(e) { this.$tools.value = e.detail }
}
document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('wall-tools', Tools)
    window.customElements.define('wall-editor', Editor)
    window.customElements.define('wall-app', Application)
})
