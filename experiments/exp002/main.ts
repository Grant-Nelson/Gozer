// The following blocks are trying to represent the following code:
//--[main.go]-------------------------------
//  package main
//  
//  func main() {
//  	println(`Starting`)
//  	for i := 0; i < 10; i++ {
//  		print(i, `. `)
//  		for j := 0; j < i; j++ {
//  			print(`*`)
//  		}
//  		println()
//  	}
//  	println(`Finished`)
//  }
//
//--[output]--------------------------------
//  Starting
//  0. 
//  1. *
//  2. **
//  3. ***
//  4. ****
//  5. *****
//  6. ******
//  7. *******
//  8. ********
//  9. *********
//  Finished
//
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

const s = new Scheduler.Scheduler();
//--[Block 0]-------------------------------
//  println(`Starting`)
//  for i := 0; ...
s.addBlock((_: any[]) => {
    print('Starting\n');
    var i = 0;
    return Scheduler.Goto(1, [i]);
})
//--[Block 1]-------------------------------
//  ... i < 10; ... {
//      print(i, `. `)
//      for j := 0;
s.addBlock((args: any[]) => {
    var i = Number(args[0]);
    if (i >= 10) return Scheduler.Goto(4, []);
    print(i + '. ');
    var j = 0;
    return Scheduler.Goto(2, [i, j]);
})
//--[Block 2]-------------------------------
//        ... j < i; j++ {
//              print(`*`)
//          }
s.addBlock((args: any[]) => {
    const i = Number(args[0]);
    var j = Number(args[1]);
    if (j >= i) return Scheduler.Goto(3, [i]);
    print('*');
    j++;
    return Scheduler.Goto(2, [i, j]);
})
//--[Block 3]-------------------------------
//   ... i++ { ...
//          println()
//      }
s.addBlock((args: any[]) => {
    var i = Number(args[0]);
    print('\n');
    i++;
    return Scheduler.Goto(1, [i]);
})
//--[Block 4]-------------------------------
//  	println(`Finished`)
//  }
s.addBlock((_: any[]) => {
    print('Finished\n');
    return Scheduler.Ret();
})

async function startMain() {
    // Start main function thread
    s.addThread(0);

    print('schedular started\n');
    while (s.isRunning()) {
        await Scheduler.Sleep(10);
    }
    print('scheduler stopped\n');
}

startMain();
