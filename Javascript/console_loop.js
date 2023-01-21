readLine = require('readline')
const { exit } = require('process');
const { promisify } = require('util')
const execPromise = promisify(require('child_process').exec);

const rl = readLine.createInterface({
    input: process.stdin,
    output: process.stdout
})

// Bind keypress event for CTRL-P interrupt
readLine.emitKeypressEvents(process.stdin)
process.stdin.setRawMode(true)
process.stdin.on('keypress', (_str, key) => {
    if (key.ctrl && key.name == 'p') {
        console.log()
        exit(0)
    }
})

const askCommand = () => {
    return new Promise(resolve => { rl.question('sh> ', resolve) })
}

const execCallback = (commandResult) => {
    console.log(`stdout: ${commandResult.stdout}`)
    loopREPL()
}

const loopREPL = () => {
    askCommand().then(execPromise).then(execCallback)
}

loopREPL()
