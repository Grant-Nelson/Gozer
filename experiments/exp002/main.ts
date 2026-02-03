// The following blocks are trying to represent the following code:
//--[main.go]-------------------------------
//  package main
//  
//  func fib(n int) int {
//  	if n <= 1 {
//  		return n;
//  	}
//  	return fib(n - 1) + fib(n - 2);
//  }
//
//  func main() {
//      println(`fib(-1) = `, fib(-1))
//      println(`fib(2) = `, fib(2))
//  	println(`fib(5) = `, fib(5))
//  	println(`fib(10) = `, fib(10))
//  }
//--[output]--------------------------------
// fib(-1) = -1
// fib(2) = 1
// fib(5) = 5
// fib(10) = 55
//------------------------------------------

import { Scheduler } from './scheduler.ts';

const divElem = document.querySelector('#root') as HTMLDivElement | null;
if (divElem != null) {
    divElem.textContent       = '';
    divElem.style.fontFamily  = 'monospace';
    divElem.style.whiteSpace  = 'pre';
    divElem.style.borderStyle = 'solid';
    divElem.style.borderColor = 'gray';
    divElem.style.padding     = '10px';
    divElem.style.overflow    = 'auto';
}

function print(text: String) {
    if (divElem == null) console.log(text);
    else divElem.textContent += text;
}

const blocks : Scheduler.Block[] = [];
blocks.push(...[
    (args: any[]) => { // Block 0 [fib]
        const n = Number(args[0]);
        if (n <= 1) return new Scheduler.Ret([n]);
        return new Scheduler.Call(blocks[1], [n], blocks[0], [n-1]);
    },
    (args: any[]) => { // Block 1
        const n = Number(args[0]);
        const ret1 = Number(args[1]);
        return new Scheduler.Call(blocks[2], [ret1], blocks[0], [n-2]);
    },
    (args: any[]) => { // Block 2
        const ret1 = Number(args[0]);
        const ret2 = Number(args[1]);
        return new Scheduler.Ret([ret1 + ret2]);
    },
    (args: any[]) => { // Block 3 [main]
        return new Scheduler.Call(blocks[4], [], blocks[0], [-1]);
    },
    (args: any[]) => { // Block 4
        print('fib(-1) = '+Number(args[0])+'\n');
        return new Scheduler.Call(blocks[5], [], blocks[0], [2]);
    },
    (args: any[]) => { // Block 5
        print('fib(2) = '+Number(args[0])+'\n');
        return new Scheduler.Call(blocks[6], [], blocks[0], [5]);
    },
    (args: any[]) => { // Block 6
        print('fib(5) = '+Number(args[0])+'\n');
        return new Scheduler.Call(blocks[7], [], blocks[0], [10]);
    },
    (args: any[]) => { // Block 7
        print('fib(10) = '+Number(args[0])+'\n');
        return new Scheduler.Ret();
    }
]);

Scheduler.Schedular.addThread(blocks[3], []);
