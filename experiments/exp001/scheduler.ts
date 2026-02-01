export namespace Scheduler {
    export interface BlockRet {
        _apply(scheduler: Scheduler, thread: Thread): void;
    }

    class GotoRet implements BlockRet {
        readonly followBlock: number;
        readonly followArgs:  any[];

        constructor(followBlock: number, followArgs: any[] = []) {
            this.followBlock = followBlock;
            this.followArgs  = followArgs; 
        }

        _apply(_: Scheduler, thread: Thread) {
            thread.addCall(this.followBlock, this.followArgs);
        }
    }

    export function Goto(block: number, args: any[] = []): BlockRet {
        return new GotoRet(block, args);
    }

    class RetRet implements BlockRet {
        readonly results: any[];

        constructor(results: any[] = []) {
            this.results = results;
        }

        _apply(_: Scheduler, thread: Thread): void {
            const len = thread.callStack.length;
            // NOTE: Need to deal with returns for an entry point call by
            // setting these values to a promise attached to a thread.
            if (len === 0) return;

            const node = thread.callStack[len-1];
            node.args.push(...this.results);
        }
    }

    export function Ret(results: any[] = []): BlockRet {
        return new RetRet(results);
    }

    export async function Sleep(ms: number) {
        await new Promise(r => setTimeout(r, ms));
    }

    export type BlockDelegate = (args: any[]) => BlockRet;

    export class Scheduler {
        _blocks: BlockDelegate[] = [];

        _nextId:   number   = 0;
        _running?: Thread   = undefined;
        _active:   Thread[] = [];
        _halt:     boolean  = false;

        // NOTE: Consider using more precise (nanoseconds) or faster time
        _swapTimeOut: number = 0;
        _pumpTimeOut: number = 0;
        readonly _swapMs: number;
        readonly _pumpMs: number;

        constructor(swapMs: number = 20, pumpMs: number = 50) {
            this._swapMs = swapMs;
            this._pumpMs = pumpMs;
        }

        addBlock(fn: BlockDelegate): number {
            return this._blocks.push(fn);
        }

        addThread(block: number, args: any[] = []): number {
            const id = this._nextId;
            this._nextId++;
            const thread = new Thread(id);
            thread.addCall(block, args);
            this._active.push(thread);

            if (this._running === undefined) this.restart();
            return id;
        }

        restart() { setTimeout(() => this._run()); }

        isRunning(): boolean { return !(this._running === undefined); }
        
        halt() { this._halt = true; }

        async _run() {
            if (this.isRunning()) return;
            this._prepare();
            while (true) {
                if (this._halt) {
                    this._halt = false;
                    return;
                }

                const now = Number(new Date());
                if (this._shouldPump(now)) await this._pump();
                if (this._shouldSwap(now)) this._swapThreads(now);

                // If running is set to undefined, then all threads have ended
                // and the schedular can stop running.
                const running = this._running;
                if (running === undefined) return;

                const call = running.callStack.pop();
                if (call === undefined) {
                    this._active.pop();
                    this._running = undefined;
                    continue;
                }

                const ret = this._blocks[call.block](call.args);
                ret._apply(this, running);
            }
        }

        _prepare() {
            const now = Number(new Date());
            this._swapTimeOut = now + this._swapMs;
            this._pumpTimeOut = now + this._pumpMs;
        }

        _shouldSwap(now: number): boolean {
            return this._running === undefined || this._swapTimeOut < now;
        }

        _shouldPump(now: number): boolean {
            return this._pumpTimeOut < now;
        }
        
        _swapThreads(now: number) {
            this._swapTimeOut = now + this._swapMs;
            // NOTE: May want to use a random selection.
            const thread = this._active.shift();
            if (thread === undefined) {
                this._running = undefined;
                return;
            }
            this._running = thread;
            // Put running thread back in active at the end.
            this._active.push(thread);
        }

        async _pump() {
            await Sleep(0);
            this._pumpTimeOut = Number(new Date()) + this._pumpMs;
        }
    }

    class Thread {
        readonly id: number;
        callStack:   CallNode[] = [];
        alive:       boolean    = true;

        constructor(id: number) { this.id = id; }

        addCall(block: number, args: any[] = []) {
            this.callStack.push(new CallNode(block, args));
        }
    }

    class CallNode {
        readonly block: number;
        readonly args:  any[];

        constructor(block: number, args: any[] = []) {
            this.block = block;
            this.args  = args; 
        }
    }
}
