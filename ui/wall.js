import crel from 'crelt'
import {Circle, Group, Image, Layer, Line, Rect, Stage, Text} from 'konva'
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
class ColorInput extends HTMLElement {
    static get observedAttributes() { return ['value'] }
    get value() { return this.getAttribute('value') }
    set value(v) { this.setAttribute('value', v) }
    connectedCallback() { this.appendChild(crel('input', {type: 'color', onchange: this.fireValueEvent, value: this.value})) }
    fireValueEvent(e) { this.dispatchEvent(new CustomEvent('value', {detail: e.target.value, bubbles: true}))}}
class FontInput extends HTMLElement {connectedCallback() { this.appendChild(crel('input', {placeholder: 'Font'})) }}
class SizeInput extends HTMLElement {
    static get observedAttributes() { return ['value'] }
    get value() { return this.getAttribute('value') }
    set value(v) { this.setAttribute('value', v) }
    connectedCallback() { this.appendChild(crel('input', {type: 'number', oninput: this.fireValueEvent, value: this.value})) }
    fireValueEvent(e) { this.dispatchEvent(new CustomEvent('value', {detail: e.target.value, bubbles: true}))}}
class DecorationsInput extends HTMLElement {connectedCallback() { this.appendChild(crel('input', {placeholder: 'Decorations'})) }}
// class SelectionConfig extends HTMLElement {}
class PenConfig extends HTMLElement {
    constructor() {
        super();
        this.color = '#ff0000'
        this.size = 16
        this.$color = crel('color-input', {value: this.color})
        this.$size = crel('size-input', {value: this.size})}
    connectedCallback() {
        this.$color.addEventListener('value', e => this.color = e.target.value)
        this.$size.addEventListener('value', e => this.size = e.target.value)
        this.appendChild(this.$size)
        this.appendChild(this.$color)}}
class TextConfig extends HTMLElement {
    connectedCallback() {
        this.appendChild(crel('font-input'))
        this.appendChild(crel('size-input'))
        this.appendChild(crel('decorations-input'))
        this.appendChild(crel('color-input'))}}
class ImageConfig extends HTMLElement {}
class ShapeConfig extends HTMLElement {}
function createTextarea(n, areaPosition) {
    console.log('text: ' + n.text())
    const textarea = crel('textarea', {class: "wall", value: n.text()});
    textarea.value = n.text()
    const style = {top: areaPosition.y+'px', left: areaPosition.x+'px', width: n.width()-n.padding()*2+'px', height: n.height()-n.padding()*2+5+'px', fontSize: n.fontSize()+'px', lineHeight: n.lineHeight(), fontFamily: n.fontFamily(), textAlign: n.align(), color: n.fill()}
    Object.keys(style).forEach(k => textarea.style[k] = style[k])
    return textarea
}
class ConfigMenu extends HTMLElement {
    static get observedAttributes() { return ['tool'] }
    get tool() { return this.getAttribute('tool') }
    set tool(v) { this.setAttribute('tool', v) }
    attributeChangedCallback(name, oldValue, newValue) {if (name === 'tool') { console.log('change tool options to: ' + newValue)}}
    constructor() {
        super(); this.$pen = crel('pen-config')
    }
    connectedCallback() {
        this.appendChild(this.$pen)
        // this.appendChild(crel('text-config'))
        // this.appendChild(crel('image-config'))
        // this.appendChild(crel('shape-config'))
    }
}
class Editor extends HTMLElement {
    static get observedAttributes() { return ['mode'] }
    get mode() { return this.getAttribute('mode') }
    set mode(v) { this.setAttribute('mode', v) }
    constructor() {
        super()
        this.isPaint = false
        this.lastPointerPosition = null
        this.scaleBy = 1.1
        this.scale = 1
        this.stageWidth = 900;
        this.stageHeight = 400;
        this.paperWidth = 900;
        this.paperHeight = 400;
        this.state = []
        this.history = [this.state]
        this.historyStep = 0
        this.lastLine = null
        this.configs = {}
    }
    connectedCallback() {
        this.dispatchEvent(new CustomEvent('inject-pen-config', {detail: {provider: this}, bubbles: true}))
        const container = document.querySelector('#wrapper');
        const width = container.offsetWidth, height = container.offsetHeight
        // const width = 900, height = 400
        this.stage = new Stage({container: this, width: width, height: height})
        this.layer = new Layer()
        this.group = new Group({x: (width / 2) - (this.paperWidth / 2), y: (height / 2) - (this.paperHeight / 2)});
        const paper = new Rect({width: this.paperWidth, height: this.paperHeight, stroke: "black", fill: "white"});
        this.group.add(paper);
        this.layer.add(this.group)
        this.stage.add(this.layer)
        this.stage.on('mousedown touchstart', _ => this.onMousedown())
        this.stage.on('mouseup touchend', _ => this.onMouseup())
        this.stage.on('mousemove touchmove', _ => this.onMousemove())
        container.addEventListener('wheel', e => this.onWheel(e))
        // this.group.on('wheel', e => this.onWheel(e))
        window.addEventListener('resize', _ => this.onResize())
        // this.onResize()
    }
    attributeChangedCallback(name, oldValue, newValue) {
        if (name === 'mode') {
            console.log('change editor mode: ' + newValue)
        }
    }
    getRelativePointerPosition(node) {
        const transform = node.getAbsoluteTransform().copy();
        transform.invert();
        const pos = node.getStage().getPointerPosition();
        return transform.point(pos);
    }
    getPointerPosition() { return this.getRelativePointerPosition(this.group)}
    create(state) {
        this.layer.destroyChildren();
        this.state.forEach((item,key) => {
        })
    }
    createShape(pos) {
        const shape = new Circle({x: pos.x, y: pos.y, fill: 'red', radius: 20})
        this.group.add(shape)
        this.layer.batchDraw()
    }
    createText(pos) {
        const fontSize = 50
        const text = new Text({text: 'hello', x: pos.x, y: pos.y-(fontSize/2), fill: 'black', fontSize: 50})
        this.group.add(text)
        this.layer.batchDraw()
    }
    createImage(pos) {
        Image.fromURL('/static/img/yoda.jpg', image => {
            this.group.add(image)
            image.position({x: pos.x, y: pos.y-(image.height()/2)})
            image.draggable(true);
            this.layer.batchDraw();
        })
    }
    createLine(pos) {
        this.isPaint = true;
        this.lastLine = new Line({stroke: this.configs.$pen.color, strokeWidth: this.configs.$pen.size, points: [pos.x, pos.y],
            globalCompositeOperation: this.mode === 'pen' ? 'source-over' : 'destination-out'})
        this.group.add(this.lastLine)
    }
    onResize() {
        const container = document.querySelector('#wrapper');
        console.log('available width: ' + container.scrollWidth)
        var containerWidth = container.offsetWidth;
        var containerHeight = container.offsetHeight;
        var scale = containerWidth / this.stageWidth;
        this.stage.width(this.stageWidth * scale );
        this.stage.height(this.stageHeight * scale);
        this.stage.scale({ x: scale, y: scale });
        this.stage.draw();
    }
    onMousedown() {
        const pos = this.getPointerPosition()
        if (this.mode === 'shape') this.createShape(pos)
        else if (this.mode === 'text') this.createText(pos)
        else if (this.mode === 'image') this.createImage(pos)
        else if (this.mode === 'pen') this.createLine(pos)
    }
    onMouseup() {
        this.isPaint = false
    }
    onMousemove() {
        if (!this.isPaint) return
        const pos = this.getPointerPosition()
        const newPoints = this.lastLine.points().concat([pos.x, pos.y]);
        this.lastLine.points(newPoints);
        this.layer.batchDraw();
    }
    onWheel(e) {
        const oldScale = this.stage.scaleX();
        const pointer = this.stage.getPointerPosition();
        const mousePointTo = {x: (pointer.x - this.stage.x()) / oldScale, y: (pointer.y - this.stage.y()) / oldScale};
        const newScale = e.deltaY > 0 ? oldScale * this.scaleBy : oldScale / this.scaleBy;
        this.stage.scale({ x: newScale, y: newScale });
        const newPos = {x: pointer.x - mousePointTo.x * newScale, y: pointer.y - mousePointTo.y * newScale};
        this.stage.position(newPos);
        this.stage.batchDraw();
    }
    onTextClick() {
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
    }
}
class App extends HTMLElement {
    constructor() {
        super()
        this.$editor = crel('wall-editor')
        this.$tools = crel('wall-tools', {})
        this.$config = crel('config-menu', {})
    }
    connectedCallback() {
        this.addEventListener('inject-pen-config', e => { e.detail.provider.configs = this.$config; e.stopPropagation() })
        this.$wallMenu = crel('div', {id: 'wall-menu'}, this.$tools, this.$config)
        this.$tools.addEventListener('value', e => this.toolValue(e))
        this.appendChild(crel('div', {id: 'wrapper'}, this.$editor))
        this.appendChild(this.$wallMenu)
    }
    toolValue(e) {
        this.$tools.value = e.detail
        this.$editor.mode = e.detail
    }
}
document.addEventListener('DOMContentLoaded', _ => {
    window.customElements.define('color-input', ColorInput)
    window.customElements.define('font-input', FontInput)
    window.customElements.define('size-input', SizeInput)
    window.customElements.define('decorations-input', DecorationsInput)
    // window.customElements.define('selection-config', SelectionConfig)
    window.customElements.define('pen-config', PenConfig)
    window.customElements.define('text-config', TextConfig)
    window.customElements.define('image-config', ImageConfig)
    window.customElements.define('shape-config', ShapeConfig)
    window.customElements.define('config-menu', ConfigMenu)
    window.customElements.define('wall-tools', Tools)
    window.customElements.define('wall-editor', Editor)
    window.customElements.define('wall-app', App)
})
