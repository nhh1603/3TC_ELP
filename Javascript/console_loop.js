readLine = require('readline')
const { promisify } = require('util')
const execPromise = promisify(require('child_process').exec);

const rl = readLine.createInterface({
    input: process.stdin,
    output: process.stdout
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