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
function createTextarea(n, areaPosition) {
    console.log('text: ' + n.text())
    const textarea = crel('textarea', {class: "wall", value: n.text()});
    textarea.value = n.text()
    const style = {top: areaPosition.y+'px', left: areaPosition.x+'px', width: n.width()-n.padding()*2+'px', height: n.height()-n.padding()*2+5+'px', fontSize: n.fontSize()+'px', lineHeight: n.lineHeight(), fontFamily: n.fontFamily(), textAlign: n.align(), color: n.fill()}
    Object.keys(style).forEach(k => textarea.style[k] = style[k])
    return textarea
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
        let circle = new Circle({x: this.stage.width() / 2, y: this.stage.height() / 2, radius: 70, fill: 'red', stroke: 'black', strokeWidth: 4})
        this.layer.add(circle)
        this.stage.add(this.layer)
        const canvas = crel('canvas', {width: this.stage.width(), height: this.stage.height()})
        this.image = new Image({image: canvas, x: 0, y: 0})
        this.layer.add(this.image)
        this.textLayer = new Layer();
        this.stage.add(this.textLayer);
        this.textNode = new Text({text: 'Some text here', x: 50, y: 80, fontSize: 20, draggable: true, width: 200})
        this.textLayer.add(this.textNode);
        this.tr = new Transformer({
            node: this.textNode, enabledAnchors: ['middle-left', 'middle-right'],
            boundBoxFunc (oldBox, newBox) {
                newBox.width = Math.max(30, newBox.width)
                return newBox}})
        this.textNode.on('transform', () => {this.textNode.setAttrs({width: this.textNode.width() * this.textNode.scaleX(), scaleX: 1});})
        this.textLayer.add(this.tr);
        this.stage.draw()
        this.context = canvas.getContext('2d')
        this.context.strokeStyle = '#df4b26';
        this.context.lineJoin = 'round';
        this.context.lineWidth = 5;
        this.textNode.on('click', () => {
            console.log('dbclick textNode')
            let n = this.textNode
            this.textNode.hide(); this.tr.hide(); this.textLayer.draw();
            // create textarea over canvas with absolute position first we need to find position for textarea how to find it? at first lets find position of text node relative to the stage
            let textPosition = n.absolutePosition();
            // then lets find position of stage container on the page:
            let stageBox = this.stage.container().getBoundingClientRect();
            // so position of textarea will be the sum of positions above:
            let areaPosition = {x: stageBox.left + textPosition.x, y: stageBox.top + textPosition.y};
            let textarea = createTextarea(n, areaPosition)
            document.body.appendChild(textarea);
            let rotation = n.rotation()
            let transform = '';
            if (rotation) { transform += 'rotateZ(' + rotation + 'deg)';}
            let px = 0;
            // also we need to slightly move textarea on firefox because it jumps a bit
            let isFirefox = navigator.userAgent.toLowerCase().indexOf('firefox') > -1;
            if (isFirefox) {px += 2 + Math.round(n.fontSize() / 20);}
            transform += 'translateY(-' + px + 'px)';
            textarea.style.transform = transform;
            // reset height
            textarea.style.height = 'auto';
            // after browsers resized it we can set actual value
            textarea.style.height = textarea.scrollHeight + 3 + 'px';
            textarea.focus();
            let removeTextarea = () => {
                textarea.parentNode.removeChild(textarea);
                window.removeEventListener('click', handleOutsideClick);
                n.show();
                this.tr.show();
                this.tr.forceUpdate();
                this.textLayer.draw();}
            function setTextareaWidth(newWidth) {
                if (!newWidth) {newWidth = n.placeholder.length * n.fontSize();}
                // some extra fixes on different browsers
                let isSafari = /^((?!chrome|android).)*safari/i.test(navigator.userAgent);
                let isFirefox = navigator.userAgent.toLowerCase().indexOf('firefox') > -1;
                if (isSafari || isFirefox) { newWidth = Math.ceil(newWidth);}
                let isEdge = document.documentMode || /Edge/.test(navigator.userAgent);
                if (isEdge) { newWidth += 1;}
                textarea.style.width = newWidth + 'px';
            }
            textarea.addEventListener('keydown',  (e) => {
                // hide on enter but don't hide on shift + enter
                if (e.key === 'Enter' && !e.shiftKey) {
                    n.text(textarea.value);
                    removeTextarea()
                } else if (e.key === 'Esc') { removeTextarea() }
            });
            textarea.addEventListener('keydown', function (e) {
                let scale = n.getAbsoluteScale().x;
                setTextareaWidth(n.width() * scale);
                textarea.style.height = 'auto';
                textarea.style.height = textarea.scrollHeight + n.fontSize() + 'px';
            });
            function handleOutsideClick(e) {
                if (e.target !== textarea) {
                    n.text(textarea.value);
                    removeTextarea();
                }
            }
            setTimeout(() => { window.addEventListener('click', handleOutsideClick); });
        });
        // this.stage.on('mousedown touchstart', _ => this.onMousedown())
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
