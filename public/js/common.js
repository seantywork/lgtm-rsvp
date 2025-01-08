USER_LOGIN = {

    id: "",
    passphrase: ""
}

async function userSignin(){



    let u_id = document.getElementById("user-id").value 

    if(u_id == ""){
  
        alert("no user id")
    
        return
    
    }


    let u_pw = document.getElementById("user-pw").value 


    if(u_pw == ""){
  
        alert("no user pw")
    
        return
    
    }



    let uc = JSON.parse(JSON.stringify(USER_LOGIN))

    uc.id = u_id
    uc.passphrase = u_pw

    let req = {
        data: JSON.stringify(uc)
    }

    let resp = await fetch(`/api/signin`, {
        body: JSON.stringify(req),
        method: "POST"
    })


    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to login")

        return
    }

    alert("successfully logged in: " + result.reply)

    location.href = "/"


}




async function getArticleList(){


      let resp = await fetch("/api/story/list", {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to get sample list")

        return

    }


    let contentReader = document.getElementById("article-reader")

    let contentEntry = JSON.parse(result.reply)
    
    console.log(contentEntry)
  
}
  
