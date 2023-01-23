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

const askCmd = () => {
    return new Promise(resolve => { rl.question('sh> ', resolve) })
}

const execCallback = (cmdResult) => {
    console.log(`${cmdResult.stdout}`)
    loopREPL()
}

const getSignal = arg => {
    if (arg === "-k") {
        return "-KILL"
    } else if (arg === "-p") {
        return "-STOP"
    } else if (arg === "-c") {
        return "-CONT"
    } else {
        signalErr = new Error()
        signalErr.stderr = "Unrecognized signal: " + arg
        throw signalErr
    }
}

const translateCmd = (cmdStr) => {
    cmdElements = cmdStr.split(" ")
    cmd = cmdElements[0] // First element is the command itself, subsequent elements are the args
    switch (cmd) {
        case "lp": // Show all open processes
            cmd = "ps -e"
            break;
        case "bing": // Change process state
            cmd = "kill"
            signalArg = cmdElements[1]
            cmdElements[1] = getSignal(signalArg)
        default:
            break;
    }

    if (cmdStr.endsWith("!")) { // Run program in the background
        // cmd = cmdStr.substring(0, cmdStr.length-1) +
        loopREPL();
        return cmdStr.substring(0, cmdStr.length-1) + "&"
    }

    if (cmd.startsWith("keep") && !isNaN(cmdElements[1])) { // Detach processus
        cmd = "disown -h %" + `{cmdElements[1]}`
    }

    cmdElements[0] = cmd
    return cmdElements.join(" ")
}

const loopREPL = () => {
    askCmd()
        .then((cmdStr) => { cmd = translateCmd(cmdStr); return execPromise(cmd) })
        .then(execCallback)
        .catch(err => {
            console.log(err.stderr)
            loopREPL()
        })
}

loopREPL()
