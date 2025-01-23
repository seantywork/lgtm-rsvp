
async function getArticleContent(){

  let contentId = CONTENT_ID

  if (contentId == "") {

      alert("no content id provided")

      return
  }

  let resp = await fetch(`/api/story/download/${contentId}`, {
      method: "GET"
  })

  let result = await resp.json()

  if(result.status != "success"){


    alert("failed to get content")

    return

  }

  let result_data = JSON.parse(result.reply)

  let blen = result_data.blocks.length

  let title = ""
  let titleIdx = -1

  let introIdx = -1
  let dateMarkedIdx = -1

  for(let i = 0; i < blen; i++){

    let b = result_data.blocks[i]

    if(b.type == "header"){
      if(b.data.level == 1){
        if(titleIdx == -1){
          title = b.data.text
          titleIdx = i
        }
      }

      if(b.data.level == 2){
        if(introIdx == -1){
          introIdx = i
        }

      }

      if(b.data.level == 3){
        if(dateMarkedIdx == -1){
          dateMarkedIdx = i
        }
      }
    }

  }

  result_data.blocks.splice(dateMarkedIdx,1)
  result_data.blocks.splice(introIdx, 1)
  result_data.blocks.splice(titleIdx, 1)
  
  let titleEl = document.getElementById("story-title")

  titleEl.innerText = title

  editor_content = new EditorJS({

    readOnly: true,
  
    holder: 'story-reader',

    tools: {

      header: {
        class: Header,
        inlineToolbar: ['marker', 'link'],
        config: {
          placeholder: 'Header'
        },
        shortcut: 'CMD+SHIFT+H'
      },


      image: {
        class: ImageTool,
      },

      list: {
        class: EditorjsList,
        inlineToolbar: true,
        shortcut: 'CMD+SHIFT+L'
      },

      checklist: {
        class: Checklist,
        inlineToolbar: true,
      },

      quote: {
        class: Quote,
        inlineToolbar: true,
        config: {
          quotePlaceholder: 'Enter a quote',
          captionPlaceholder: 'Quote\'s author',
        },
        shortcut: 'CMD+SHIFT+O'
      },

      warning: Warning,

      marker: {
        class:  Marker,
        shortcut: 'CMD+SHIFT+M'
      },

      code: {
        class:  CodeTool,
        shortcut: 'CMD+SHIFT+C'
      },

      delimiter: Delimiter,

      inlineCode: {
        class: InlineCode,
        shortcut: 'CMD+SHIFT+C'
      },

      linkTool: LinkTool,

      embed: Embed,

      table: {
        class: Table,
        inlineToolbar: true,
        shortcut: 'CMD+ALT+T'
      },

    },


    data: {
      blocks: result_data.blocks
    },
    onReady: async function(){

      console.log("data ready")

    },
//      onChange: function(api, event) {
//        console.log('something changed', event);
//      }
  })


}




(async function(){

  await getArticleContent()

})()

