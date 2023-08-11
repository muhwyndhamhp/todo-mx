var addTodoPath = "/add-todo"

var simplemde = new SimpleMDE({ 
    element: document.getElementById("todo-body"),
    autosave: {
        enable: false
    },
    forceSync: true,
});

document.body.addEventListener('htmx:configRequest', function(evt){
    if (evt.detail.path != addTodoPath ){
        return
    }
    params = evt.detail.parameters
    html = simplemde.options.previewRender(params['body'])
    evt.detail.parameters['body_encoded'] = html
})

document.body.addEventListener('htmx:afterRequest', function(evt) {
    if (evt.detail.pathInfo.finalRequestPath == addTodoPath ){
        simplemde.value("")
    }
})