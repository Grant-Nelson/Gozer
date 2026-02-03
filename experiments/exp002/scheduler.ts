export namespace Scheduler {
    export interface BlockResult {
        _apply(thread: Thread): void;
    }

    export class Goto implements BlockResult {
        readonly followBlock: Block;
        readonly followArgs:  any[];

        constructor(followBlock: Block, followArgs: any[] = []) {
            this.followBlock = followBlock;
            this.followArgs  = followArgs; 
        }

        _apply(thread: Thread) {
            thread.addCall(this.followBlock, this.followArgs);
        }
    }

    export class Ret implements BlockResult {
        readonly results: any[];

        constructor(results: any[] = []) {
            this.results = results;
        }

        _apply(thread: Thread): void {
            const len = thread.callStack.length;
            // NOTE: Need to deal with returns for an entry point call by
            // setting these values to a promise attached to a thread.
            if (len === 0) return;

            // NOTE: Needs to skip over defer blocks and any blocks before
            // the follow after a call block.
            const node = thread.callStack[len-1];
            node.args.push(...this.results);
        }
    }

    export class Call implements BlockResult {
        readonly followBlock: Block;
        readonly followArgs:  any[];
        readonly callBlock:   Block;
        readonly callArgs:    any[];

        constructor(followBlock: Block, followArgs: any[], callBlock: Block, callArgs: any[]) {
            this.followBlock = followBlock;
            this.followArgs  = followArgs;
            this.callBlock   = callBlock;
            this.callArgs    = callArgs; 
        }

        _apply(thread: Thread): void {
            thread.addCall(this.followBlock, this.followArgs);
            thread.addCall(this.callBlock, this.callArgs);
        }
    }

    export async function Sleep(ms: number) {
        await new Promise(r => setTimeout(r, ms));
    }

    export type Block = (args: any[]) => BlockResult;

    class Thread {
        readonly id: number;
        callStack:   CallNode[] = [];

        constructor(id: number) { this.id = id; }

        addCall(block: Block, args: any[] = []) {
            this.callStack.push(new CallNode(block, args));
        }
    }

    class CallNode {
        readonly block: Block;
        readonly args:  any[];

        constructor(block: Block, args: any[] = []) {
            this.block = block;
            this.args  = args; 
        }
    }

    enum status {
        Starting,
        Started,
        Stopping,
        Stopped,
    }

    class scheduler {
        _nextThreadId: number   = 0;
        _running?:     Thread   = undefined;
        _active:       Thread[] = [];
        _status:       status   = status.Stopped;

        // NOTE: Consider using more precise (nanoseconds) or faster time
        _swapTimeOut: number = 0;
        _pumpTimeOut: number = 0;
        readonly _swapMs: number;
        readonly _pumpMs: number;

        constructor(swapMs: number = 20, pumpMs: number = 50) {
            this._swapMs = swapMs;
            this._pumpMs = pumpMs;
        }

        addThread(block: Block, args: any[] = []): number {
            const id = this._nextThreadId;
            this._nextThreadId++;
            const thread = new Thread(id);
            thread.addCall(block, args);
            this._active.push(thread);

            if (!this.isRunning()) this.restart();
            return id;
        }

        restart() {
            if (this._status === status.Stopped || this._status === status.Stopping) {
                this._status = status.Starting;
                setTimeout(() => this._run());
            }
        }

        isRunning(): boolean {
            return this._status === status.Starting || this._status === status.Started;
        }
        
        halt() {
            if (this.isRunning()) this._status = status.Stopping;
        }

        async _run() {
            if (this._status !== status.Starting) return;
            this._prepare();
            while (true) {
                // @ts-ignore: External events may change _status so warning is wrong.
                if (this._status === status.Stopping) break;

                const now = Number(new Date());
                if (this._shouldPump(now)) await this._pump();
                if (this._shouldSwap(now)) this._swapThreads(now);

                // If running is set to undefined, then all threads have ended
                // and the schedular can stop running.
                const running = this._running;
                if (running === undefined) break;

                // Get current running block for the thread.
                const call = running.callStack.pop();
                if (call === undefined) {
                    // If there is no block then the thread has exited.
                    // The BlockRet at the end of the prior block should have
                    // handled this, so this is only to handle an unexpected state.
                    this._active.pop();
                    this._running = undefined;
                    continue;
                }

                // Run the block and get return.
                // NOTE: This could be optimized to keep thread cycling without
                //   checking as much until a swap is needed.
                const ret = call.block(call.args);
                ret._apply(running);
            }
            this._shutdown();
        }

        _prepare() {
            const now = Number(new Date());
            this._swapTimeOut = now + this._swapMs;
            this._pumpTimeOut = now + this._pumpMs;
            this._status = status.Started;
        }

        _shutdown() {
            this._status = status.Stopped;
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
    
    export const Schedular = new scheduler();
}
