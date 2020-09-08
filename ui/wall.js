import crel from 'crelt'
import konva from 'konva'
class Tools extends HTMLElement {
    static get observedAttributes() { return ['value'] }
    get value() { return this.getAttribute('value') }
    set value(v) { this.setAttribute('value', v) }
    constructor() { super(); this.tools = ['selection', 'pen', 'text', 'image'] }
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
        this.appendChild(crel('h1', {}, 'Pen'))
        this.appendChild(crel('label', {'for': 'strokeSize'}, 'Stroke size')); this.appendChild(this.$size)
        this.appendChild(crel('label', {'for': 'strokeColor'}, 'Stroke color')); this.appendChild(this.$color)}}
class TextConfig extends HTMLElement {
    constructor() {
        super();
        this.$color = crel('color-input', {value: this.color})
        this.$size = crel('size-input', {value: this.size})}
    connectedCallback() {
        this.appendChild(crel('h1', {}, 'Text'))
        this.appendChild(crel('label', {'for': 'fontSize'}, 'Font size')); this.appendChild(this.$size)
        this.appendChild(crel('label', {'for': 'textColor'}, 'Text color')); this.appendChild(this.$color)}}
class ImageConfig extends HTMLElement {
    connectedCallback() {
        this.appendChild(crel('h1', {}, 'Image'))}
}
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
    attributeChangedCallback(name, oldValue, newValue) {
        if (name === 'tool') {
            this.display(oldValue, 'none')
            this.display(newValue, 'flex')
        }
    }
    display(name, value) {
        if (name === 'pen') this.$pen.style.display = value
        else if (name === 'text') this.$text.style.display = value
        else if (name === 'image') this.$image.style.display = value
    }
    constructor() {
        super();
        this.$pen = crel('pen-config')
        this.$text = crel('text-config')
        this.$image = crel('image-config')
    }
    connectedCallback() {
        this.appendChild(this.$pen)
        this.appendChild(this.$text)
        this.appendChild(this.$image)
    }
}
class TopMenu extends HTMLElement {
    connectedCallback() {
        this.appendChild(crel('div', {class: 'delete', onclick: _ => this.dispatchEvent(new CustomEvent('delete', {bubbles: true}))}))
        this.appendChild(crel('div', {class: 'undo', onclick: _ => this.dispatchEvent(new CustomEvent('undo', {bubbles: true}))}))
        this.appendChild(crel('div', {class: 'redo', onclick: _ => this.dispatchEvent(new CustomEvent('redo', {bubbles: true}))}))
    }
}
export class EditorState {
    get image() { return this.state.image }
    get text() { return this.state.text }
    get pen() { return this.state.pen }
    constructor() {
        // image / text / pen
        this.state = {image:[], text: [], pen: []}
        this.undoStack = []
        this.redoStack = []
    }
    layerName(className) {
        if (className === 'Image') return 'image'
        else if (className === 'Text') return 'text'
        else if (className === 'Line') return 'pen'
    }
    add(node) {
        this.save()
        const attrs = node.getAttrs()
        const className = node.getClassName()
        const layerName = this.layerName(className)
        this.state[layerName].push({className, attrs})
    }
    remove(node) {
        this.save()
        const className = node.getClassName()
        const layerName = this.layerName(className)
        this.state[layerName] = this.state[layerName].filter(e => e.attrs.id !== node.attrs.id)
    }
    update(node) {
        this.save()
        console.log(node)
        const className = node.getClassName()
        const layerName = this.layerName(className)
        const i = this.state[layerName].findIndex(e => e.attrs.id === node.attrs.id)
        console.log(this.state[i])
        this.state[layerName][i].attrs.x = node.x()
        this.state[layerName][i].attrs.y = node.y()
    }
    save() {
        this.undoStack.push(JSON.stringify(this.state))
        this.redoStack = []
    }
    undo() {
        if (this.undoStack.length === 0) return
        this.redoStack.push(JSON.stringify(this.state))
        this.state = JSON.parse(this.undoStack.pop())
    }
    redo() {
        if (this.redoStack.length === 0) return
        this.undoStack.push(JSON.stringify(this.state))
        this.state = JSON.parse(this.redoStack.pop())
    }
}
function toNode(el) {
    const type = el.className
    if (type === 'Line') return new konva.Line(el.attrs)
    else if (type === 'Text') return new konva.Text(el.attrs)
    else if (type === 'Image') return new konva.Image(el.attrs)
    else if (type === 'Circle') return new konva.Circle(el.attrs)
    else console.log('WTF: ' + type)
}
window.ED = null
function randomID() { return '_' + Math.random().toString(36).substr(2, 9) }
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
        this.lastLine = null
        this.configs = {}
        this.state = new EditorState()
        window.ED = this
    }
    connectedCallback() {
        this.dispatchEvent(new CustomEvent('inject-pen-config', {detail: {provider: this}, bubbles: true}))
        const container = document.querySelector('#wrapper');
        const width = container.offsetWidth, height = container.offsetHeight
        this.stage = new konva.Stage({container: this, width: width, height: height})
        this.imageLayer = new konva.Layer()
        this.textLayer = new konva.Layer()
        this.penLayer = new konva.Layer()
        this.stage.add(this.imageLayer, this.textLayer, this.penLayer)
        this.redraw()
        this.stage.on('mousedown touchstart', async e => await this.onMousedown(e))
        this.stage.on('mouseup touchend', _ => this.onMouseup())
        this.stage.on('mousemove touchmove', _ => this.onMousemove())
        this.stage.on('dragend', e => this.onDragend(e))
        container.addEventListener('wheel', e => this.onWheel(e))
    }
    attributeChangedCallback(name, oldValue, newValue) {
        if (name === 'mode') {}
    }
    redraw() {
        console.log('redraw')
        this.imageLayer.destroyChildren();
        this.textLayer.destroyChildren();
        this.penLayer.destroyChildren();
        this.state.image.forEach(e => this.imageLayer.add(toNode(e)))
        this.state.text.forEach(e => this.textLayer.add(toNode(e)))
        this.state.pen.forEach(e => this.penLayer.add(toNode(e)))
        this.stage.draw()
        // this.stage.draw();
    }
    createText(pos) {
        const id = randomID()
        const fontSize = 50
        const options = {id: id, text: 'hello', x: pos.x, y: pos.y-(fontSize/2), fill: 'black', fontSize: 50}
        const text = new konva.Text(options)
        text.draggable(true);
        this.textLayer.add(text)
        this.textLayer.batchDraw()
        this.state.add(text)
    }
    async createImage(pos) {
        const id = randomID()
        let image = await this.loadImage('/static/img/yoda.jpg')
        const options = {id: id, x: pos.x, y: pos.y-(image.height()/2), src: '/static/img/yoda.jpg'}
        this.imageLayer.add(image)
        image.position(options)
        image.id(id)
        image.draggable(true);
        this.imageLayer.batchDraw();
        this.state.add(image)
    }
    createStroke(pos) {
        const id = randomID()
        this.isPaint = true;
        const options = {id: id, stroke: this.configs.$pen.color, strokeWidth: this.configs.$pen.size, points: [pos.x, pos.y],
            globalCompositeOperation: this.mode === 'pen' ? 'source-over' : 'destination-out'}
        this.lastLine = new konva.Line(options)
        this.penLayer.add(this.lastLine)
        this.state.add(this.lastLine)
    }
    undo() {
        this.state.undo()
        this.redraw()
    }
    redo() {
        this.state.redo()
        this.redraw()
    }
    delete() {
        if (this.mode !== 'selection' || !this.selected) return
        this.state.remove(this.selected)
        this.selected.destroy()
        this.selected = null
        this.redraw()
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
    async onMousedown(e) {
        const pos = this.stage.getPointerPosition()
        if (this.mode === 'selection') {
            if (e.target.getClassName === 'Stage') return
            this.selected = e.target
        }
        else if (this.mode === 'text')  this.createText(pos)
        else if (this.mode === 'image') await this.createImage(pos)
        else if (this.mode === 'pen') this.createStroke(pos)
        else return
    }
    onMouseup() {
        this.isPaint = false
    }
    onMousemove() {
        if (!this.isPaint) return
        const pos = this.stage.getPointerPosition()
        const newPoints = this.lastLine.points().concat([pos.x, pos.y]);
        this.lastLine.points(newPoints);
        this.penLayer.batchDraw();
    }
    onDragend(e) {
        this.state.update(e.target)
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
            this.textLayer.draw();
        }
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
    loadImage(url) {
        return new Promise((resolve, reject) => {
            konva.Image.fromURL(url, image => resolve(image))
        })
    }
}
class App extends HTMLElement {
    constructor() {
        super()
        this.$editor = crel('wall-editor')
        this.$tools = crel('wall-tools', {})
        this.$config = crel('config-menu', {})
        this.$topMenu = crel('top-menu', {})
    }
    connectedCallback() {
        this.addEventListener('inject-pen-config', e => { e.detail.provider.configs = this.$config; e.stopPropagation() })
        this.$wallMenu = crel('div', {id: 'wall-menu'}, this.$tools, this.$config)
        this.$tools.addEventListener('value', e => this.toolValue(e))
        this.addEventListener('delete', _ => this.$editor.delete())
        this.addEventListener('undo', _ => this.$editor.undo())
        this.addEventListener('redo', _ => this.$editor.redo())
        this.appendChild(crel('div', {id: 'wrapper'}, this.$editor))
        this.appendChild(this.$wallMenu)
        this.appendChild(this.$topMenu)
    }
    toolValue(e) {
        this.$tools.value = e.detail
        this.$editor.mode = e.detail
        this.$config.tool = e.detail
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
    window.customElements.define('config-menu', ConfigMenu)
    window.customElements.define('top-menu', TopMenu)
    window.customElements.define('wall-tools', Tools)
    window.customElements.define('wall-editor', Editor)
    window.customElements.define('wall-app', App)
})
