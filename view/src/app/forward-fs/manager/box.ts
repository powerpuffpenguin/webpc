export class Point {
    constructor(public x: number = 0, public y: number = 0) {
    }
    toView(): Point {
        if (document.compatMode == "BackCompat") {
            this.x -= document.body.scrollLeft
            this.y -= document.body.scrollTop
        } else {
            this.x -= document.documentElement.scrollLeft
            this.y -= document.documentElement.scrollTop
        }
        return this
    }
}
// 有效範圍
export class Box {
    private p0_ = new Point()
    private p1_ = new Point()
    start = new Point()
    stop = new Point()
    setRange(element: any) {
        this.p0_ = getViewPoint(element)
        this.p1_ = new Point(this.p0_.x + element.offsetWidth, this.p0_.y + element.offsetHeight)
    }
    private _fixStart() {
        if (this.start.x < this.p0_.x) {
            this.start.x = this.p0_.x
        } else if (this.start.x > this.p1_.x) {
            this.start.x = this.p1_.x
        }

        if (this.start.y < this.p0_.y) {
            this.start.y = this.p0_.y
        } else if (this.start.y > this.p1_.y) {
            this.start.y = this.p1_.y
        }
    }

    private _fixStop() {
        if (this.stop.x < this.p0_.x) {
            this.stop.x = this.p0_.x
        } else if (this.stop.x > this.p1_.x) {
            this.stop.x = this.p1_.x
        }

        if (this.stop.y < this.p0_.y) {
            this.stop.y = this.p0_.y
        } else if (this.stop.y > this.p1_.y) {
            this.stop.y = this.p1_.y
        }
    }
    calculate() {
        if (!this.start || !this.stop) {
            return
        }
        if (this.p0_ && this.p1_) {
            this._fixStart()
            this._fixStop()
        }
        this.x = Math.min(this.start.x, this.stop.x)
        this.y = Math.min(this.start.y, this.stop.y)
        this.w = Math.abs(this.start.x - this.stop.x)
        this.h = Math.abs(this.start.y - this.stop.y)
    }
    x = 0
    y = 0
    w = 0
    h = 0
    reset() {
        this.x = 0
        this.y = 0
        this.w = 0
        this.h = 0
        this.p0_ = new Point()
        this.p1_ = new Point()
        this.start = new Point()
        this.stop = new Point()
    }

    checked(doc: Document): Array<number> {
        const result = new Array<number>()
        const nodes = doc.childNodes

        if (nodes && nodes.length > 0) {
            let parent: any
            for (let i = 0; i < nodes.length; i++) {
                let node = (nodes[i] as any)
                if (!node || !node.querySelector) {
                    continue
                }
                node = node.querySelector('.wrapper')
                if (!node) {
                    continue
                }
                const l = getViewPoint(node)
                const r = new Point(l.x + node.offsetWidth, l.y + node.offsetHeight)
                const ok = this.testView(l, r)
                if (ok) {
                    result.push(i)
                }
            }
        }
        return result
    }
    testView(l: Point, r: Point): boolean {
        if (r.x < this.x || l.x > (this.x + this.w)) {
            return false
        }
        if (r.y < this.y || l.y > (this.y + this.h)) {
            return false
        }
        return true
    }
}
export function getPagePoint(element: any): Point {
    let x = 0
    let y = 0
    while (element) {
        x += element.offsetLeft + element.clientLeft
        y += element.offsetTop + element.clientTop
        element = element.offsetParent
    }
    return new Point(x, y)
}

export function getViewPoint(element: any): Point {
    return getPagePoint(element).toView()
}