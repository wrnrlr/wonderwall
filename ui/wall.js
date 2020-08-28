import crel from 'crelt'
import {Circle, Layer, Stage, Image, Text, Transformer} from 'konva'
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
class Editor extends HTMLElement {
    constructor() {
        super()
        this.mode = 'brush'
        this.isPaint = false
        this.lastPointerPosition = null
    }
    connectedCallback() {
        this.stage = new Stage({container: this, width: 800, height: 600})
        this.layer = new Layer()
        let circle = new Circle({
            x: this.stage.width() / 2, y: this.stage.height() / 2,
            radius: 70, fill: 'red', stroke: 'black', strokeWidth: 4})
        this.layer.add(circle)
        this.stage.add(this.layer)
        const canvas = crel('canvas', {width: this.stage.width(), height: this.stage.height()})
        this.image = new Image({image: canvas, x: 0, y: 0})
        this.layer.add(this.image)
        this.stage.draw()
        this.context = canvas.getContext('2d')
        this.context.strokeStyle = '#df4b26';
        this.context.lineJoin = 'round';
        this.context.lineWidth = 5;
        this.stage.on('mousedown touchstart', _ => this.onMousedown())
        this.stage.on('mouseup touchend', _ => this.onMouseup())
        this.stage.on('mousemove touchmove', _ => this.onMousemove())
    }
    disconnectedCallback() {}
    onMousedown() {
        this.isPaint = true
        this.lastPointerPosition = this.stage.getPointerPosition()
    }
    onMouseup() {
        this.isPaint = false
    }
    onMousemove() {
        if (!this.isPaint) return
        if (this.mode === 'brush') this.context.globalCompositeOperation = 'source-over'
        this.context.beginPath()
        let localPos = {x: this.lastPointerPosition.x - this.image.x(), y: this.lastPointerPosition.y - this.image.y()}
        this.context.moveTo(localPos.x, localPos.y)
        let pos = this.stage.getPointerPosition()
        localPos = {x: pos.x - this.image.x(), y: pos.y - this.image.y()}
        this.context.lineTo(localPos.x, localPos.y)
        this.context.closePath()
        this.context.stroke()
        this.lastPointerPosition = pos
        this.layer.batchDraw()
    }
}
class App extends HTMLElement {
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
    window.customElements.define('wall-app', App)
})
