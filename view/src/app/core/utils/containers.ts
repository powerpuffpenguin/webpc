export class Pair<T0, T1>{
    constructor(public first: T0, public second: T1) { }
}
export function makePair<T0, T1>(first: T0, second: T1) {
    return new Pair<T0, T1>(first, second)
}
export class Value<T> {
    private readonly done_: boolean
    get done(): boolean {
        return this.done_
    }
    private readonly value_?: T
    get value(): T {
        if (!this.done_) {
            throw new Error("value not done")
        }
        const result: any = this.value_
        return result
    }
    constructor(
        done = false,
        value?: T,
    ) {
        if (done) {
            this.done_ = true
            this.value_ = value
        } else {
            this.done_ = false
        }
    }
}
export const noValue = new Value<any>(false)
export function makeValue<T>(value: T) {
    return new Value(true, value)
}

export class Queue<T> {
    private size_ = 0
    private offset_ = 0
    private data_: Array<T>
    constructor(size: number) {
        size = Math.floor(size)
        if (size < 1) {
            throw new Error(`size must > 0`)
        }
        this.data_ = new Array<T>(size)
        this.size_ = 0
    }
    toArray(): Array<T> {
        const length = this.length
        const data = this.data_
        const result = new Array<T>(length)
        const offset = this.offset
        const capacity = this.capacity
        for (let i = 0; i < length; i++) {
            result[i] = data[(offset + i) % capacity]
        }
        return result
    }
    get isFull(): boolean {
        return this.size_ == this.data_.length
    }
    get isEmpty(): boolean {
        return this.size_ == 0
    }
    get length(): number {
        return this.size_
    }
    get capacity(): number {
        return this.data_.length
    }
    get data(): Array<T> {
        return this.data_
    }
    get offset(): number {
        return this.offset_
    }
    clear() {
        this.offset_ = 0
        this.size_ = 0
    }
    front(): Value<T> {
        const size = this.size_
        if (size == 0) {
            return noValue
        }
        return makeValue(this.data_[this.offset_])
    }
    back(): Value<T> {
        const size = this.size_
        if (size == 0) {
            return noValue
        }
        const length = this.data_.length
        const offset = (this.offset_ + size - 1) % length

        return makeValue(this.data_[offset])
    }
    pop(): Value<T> {
        const size = this.size_
        if (size == 0) {
            return noValue
        }
        const offset = this.offset_
        this.offset_++
        this.size_--
        if (this.offset_ == this.data_.length) {
            this.offset_ = 0
        }
        return makeValue(this.data_[offset])
    }
    push(data: T): Value<T> {
        let result: Value<T>
        const size = this.size_
        const length = this.data_.length
        const offset = (this.offset_ + size) % length
        if (size == length) {
            result = makeValue(this.data_[offset])
            this.offset_++
            if (this.offset_ == this.data_.length) {
                this.offset_ = 0
            }
        } else {
            result = noValue
            this.size_++
        }
        this.data_[offset] = data
        return result
    }

}