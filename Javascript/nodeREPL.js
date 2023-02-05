readLine = require('readline')
const { exit } = require('process');
const { spawn } = require('child_process');

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
        Object.values(children).forEach(child => {
            child.kill("SIGHUP")
        })
        exit(0)
    }
})

function askCmd() {
    return new Promise(resolve => { rl.question('sh> ', resolve) })
}


// Converts NodeJs bing command's arguments to bash parameters
function getSignalType(arg) {
    if (arg === "-k") {
        return "-KILL"
    } else if (arg === "-p") {
        return "-STOP"
    } else if (arg === "-c") {
        return "-CONT"
    } else {
        throw new Error("Unrecognized signal: " + arg)
    }
}

function translateCmd(input) {
    result = {}
    inputTokens = input.split(" ")
    cmd = inputTokens[0] // First element is the command itself, subsequent elements are the args
    inputTokens.splice(0, 1)
    cmdArgs = inputTokens
    switch (cmd) {
        case "lp": // Show all open processes
            if (cmdArgs.length > 0) {
                throw new Error("lp does not take any arguments")
            }
            result.cmd = "ps"
            result.args = ["-e"]
            break;
        case "bing": // Change a process' state
            result.cmd = "kill"
            signalArg = cmdArgs[0]
            cmdArgs[0] = getSignalType(signalArg)
            result.args = cmdArgs
            break
        case "keep":
            childPid = parseInt(cmdArgs[0])
            delete children[childPid]
            return
        default:
            result.cmd = cmd
            result.args = cmdArgs
            break
    }

    return result
}

// Used to store all child processes, with their PID as keys.
children = {}

function removeFromPending(child) {
    if (child !== undefined) {
        delete children[child.pid]
    }
}

function loopREPL() {
    askCmd()
        .then((input) => {
            if (input === "") {
                return
            }
            doBackground = input.endsWith("!")
            if (doBackground) {
                input = input.substring(0, input.length - 1)
            }
            bashCmd = translateCmd(input)
            if (bashCmd === undefined) {
                return
            }
            cmd = bashCmd.cmd
            args = bashCmd.args

            spawnedChild = spawn(cmd, args)
            children[spawnedChild.pid] = spawnedChild

            if (!doBackground) {
                spawnedChild.stdout.setEncoding("utf8")
                spawnedChild.stderr.setEncoding("utf8")
                spawnedChild.stderr.on('data', data => console.log(data))
                spawnedChild.stdout.on('data', data => console.log(data))

                return new Promise((resolve, reject) => {
                    spawnedChild.on('error', reject)
                    // In foreground mode, we resolve upon exit of the child process
                    spawnedChild.on('exit', () => resolve(spawnedChild))
                })
            } else {
                return new Promise((resolve, reject) => {
                    // In background mode, we resolve upon spawning of the child process
                    spawnedChild.on('spawn', () => {
                        console.log("[Bg]", spawnedChild.pid)
                        resolve()
                    })
                    spawnedChild.on('error', reject)
                    spawnedChild.on('exit', () => removeFromPending(spawnedChild))
                })    
            }
        })
        .then(removeFromPending)
        .catch(err => console.log("Unrecognized command: " + err.path)
)
        .finally(loopREPL)
}

loopREPL()
