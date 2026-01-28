namespace Scheduler {

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

        _apply(scheduler: Scheduler, thread: Thread) {
            scheduler._addCall(thread, this.followBlock, this.followArgs);
        }
    }

    export function Goto(block: number, args: any[] = []): BlockRet {
        return new GotoRet(block, args);
    }

    export type BlockDelegate = (args: any[]) => BlockRet;

    export class Scheduler {
        _blocks: BlockDelegate[] = [];

        _nextId:   number   = 0;
        _running?: Thread   = undefined;
        _active:   Thread[] = [];

        // TODO: Consider using more precise (nanoseconds) or faster time
        _swapTimeOut: number;
        _pumpTimeOut: number;
        readonly _swapMs: number;
        readonly _pumpMs: number;

        constructor(swapMs: number = 20, pumpMs: number = 50) {
            this._swapMs = swapMs;
            this._pumpMs = pumpMs;

            const now = Number(new Date());
            this._swapTimeOut = now + this._swapMs;
            this._pumpTimeOut = now + this._pumpMs;
        }

        addBlock(fn: BlockDelegate): number {
            return this._blocks.push(fn);
        }

        addThread(block: number, args: any[] = []): number {
            const id = this._nextId;
            this._nextId++;
            const thread = new Thread(id);
            this._addCall(thread, block, args);
            this._active.push(thread);

            if (this._running === undefined) {
                setTimeout(() => this._run());
            }

            return id;
        }
        
        _addCall(thread: Thread, block: number, args: any[]) {
            thread.callStack.push(new CallNode(block, args));
        }

        async _run() {
            while (true) {
                const now = Number(new Date());
                if (this._shouldSwap(now)) this._swapThreads(now);

                // If running is set to undefined, then all threads have ended
                // and the schedular can stop running.
                const running = this._running;
                if (running === undefined) return;

                if (this._shouldPump(now)) await this._pump(now);



                const call = running.callStack.pop();
                if (call === undefined) {
                    // TODO: deal with an empty thread.
                    return
                }
                const ret = this._blocks[call.block](call.args);
                ret._apply(this, running);o
                // TODO: finish implementing

            }
        }

        _shouldSwap(now: Number): boolean {
            if (this._running === undefined) return true;



            return false;
        }
        
        _swapThreads(now: Number) {

        }

        _shouldPump(now: Number): boolean {



            return this._pumpTimer
        }

        async _pump(now: Number) {

        }
    }

    class Thread {
        readonly id: number;
        callStack:   CallNode[] = [];

        constructor(id: number) {
            this.id = id;
        }

        isMain(): boolean { return this.id == 0; }
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
