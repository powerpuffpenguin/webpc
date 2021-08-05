import { Injectable } from "@angular/core"
import { BehaviorSubject } from "rxjs"

export const RootID = '1'
export interface NetElement {
    id: string
    name: string
    description: string
    children: Array<string>
}
export class Element {
    parent: Element | undefined
    children: Array<Element> = []
    constructor(public id: string,
        public name: string,
        public description: string,
    ) { }
    get root(): boolean {
        return this.id == RootID
    }
    get leaf(): boolean {
        return this.children.length != 0
    }
    addChild(child: Element) {
        const parent = child.parent
        if (parent == this) {
            return
        }

        this.children.push(child)
        if (parent) {
            parent.removeChild(child)
        }

        child.parent = this
    }
    removeChild(child: Element) {
        const parent = child.parent
        if (parent == this) {
            return
        }

        child.parent = undefined
        const index = this.children.indexOf(child)
        this.children.splice(index, 1)
    }
}
export class Tree {
    keys = new Map<string, Element>()
    root: Element | undefined
    constructor(public readonly helper: Helper, items: Array<NetElement>) {
        const keys = new Map<string, Element>()
        items?.forEach((item) => {
            keys.set(item.id, new Element(item.id, item.name, item.description))
        })
        const root = keys.get(RootID)
        if (!root) {
            throw new Error('root not exists')
        }
        items.forEach((item) => {
            const parent = keys.get(item.id)
            item.children?.forEach((id) => {
                const child = keys.get(id)
                if (child) {
                    child.parent = parent
                    parent?.children.push(child)
                }
            })
        })
        this.root = root
        this.forEach((ele) => {
            this.keys.set(ele.id, ele)
        })
        // update
        const rootNode = new NestedNode(root)
        this._appendChildren(rootNode)
        helper.dataChange.next([rootNode])
    }
    private _appendChildren(node: NestedNode) {
        const ele = node.data
        ele.children.forEach((child) => {
            const nodeChild = new NestedNode(child)
            nodeChild.parent = node
            node.children.push(nodeChild)
            this._appendChildren(nodeChild)
        })

        node.children.sort(NestedNode.compareFn)
    }

    forEach(callback: (ele: Element) => void, root?: Element) {
        const ele = root ?? this.root
        if (ele) {
            callback(ele)
            ele?.children?.forEach((child) => {
                this.forEach(callback, child)
            })
        }
    }
    add(node: NestedNode, data: Element) {
        const keys = this.keys
        const parent = keys.get(node.data.id)
        if (!parent) {
            throw new Error(`parent not exists: ${node.data.id}`)
        } else if (parent != node.data) {
            throw new Error(`parent expired: ${node.data.id}`)
        } else if (keys.has(data.id)) {
            throw new Error(`id already exists: ${node.data.id}`)
        }
        parent.addChild(data)
        // add view
        const child = new NestedNode(data)
        child.parent = node
        node.children.push(child)
        node.children.sort(NestedNode.compareFn)

        keys.set(data.id, data)
        // update view
        this.helper.updateView()
    }
    remove(current: Element) {
        const keys = this.keys
        if (keys.get(current.id) != current) {
            throw new Error(`parent expired: ${current.id}`)
        }
        this.forEach((ele) => {
            keys.delete(ele.id)
        }, current)
        current.parent?.removeChild(current)
    }
}


export class NestedNode {
    parent: NestedNode | undefined
    constructor(public data: Element,
    ) { }
    children: NestedNode[] = []
    static compareFn(l: NestedNode, r: NestedNode): number {
        const lv = l.data?.name ?? ''
        const rv = r.data?.name ?? ''
        if (lv > rv) {
            return 1
        } else if (lv == rv) {
            return 0
        }
        return -1
    }
}
export class FlatNode {
    constructor(public data: Element,
    ) { }
    expandable = false
    level = 0
    update = false
}
@Injectable()
export class Helper {
    dataChange = new BehaviorSubject<NestedNode[]>([])
    get data(): NestedNode[] { return this.dataChange.value }
    constructor() {
        const data = new Array<NestedNode>()
        this.dataChange.next(data)
    }
    updateView() {
        this.dataChange.next(this.data)
    }
    add(parent: NestedNode, ele: Element) {
        parent.children.push(new NestedNode(ele))
        this.dataChange.next(this.data)
    }
}