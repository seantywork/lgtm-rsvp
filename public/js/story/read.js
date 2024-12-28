
async function getArticleContent(){

  let contentId = CONTENT_ID

  if (contentId == "") {

      alert("no content id provided")

      return
  }

  let resp = await fetch(`/api/story/r/${contentId}`, {
      method: "GET"
  })

  let result = await resp.json()

  if(result.status != "success"){


    alert("failed to get content")

    return

  }

  let result_data = JSON.parse(result.reply)

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




async function loadArticleContent(){



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
      blocks: [
        {
          type : 'paragraph',
          data : {
            text : '테스트 입니다'
          }
        },
        {
          type : 'paragraph',
          data : {
            text : '테스트 입니다 테스트 입니다 테스트 입니다 테스트 입니다 테스트 입니다 테스트 입니다 테스트 입니다 '
          }
        },
        {
          type : 'paragraph',
          data : {
            text : '테스트 입니다'
          }
        }
      ]
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

  await loadArticleContent()

})()

