function send() {
    const msg = document.getElementById("msg").value

    fetch("/send", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
        body: "msg=" + msg 
    })
}

async function initQueue() {
    await fetch("/createQueue", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
    });
}

async function main(){
    await initQueue()
}

main()