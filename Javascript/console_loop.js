readLine = require('readline')

const rl = readLine.createInterface({
    input: process.stdin,
    output: process.stdout
})

const askCommand = () => {
    return new Promise(resolve => { rl.question('sh> ', resolve) } )
}

const promptProgram = async () => {
    let reply
    while (reply !== "end") {
        reply = await askCommand()
        console.log(reply)
    }
    rl.close()
}

promptProgram()