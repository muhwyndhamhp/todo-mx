var simplemde = new SimpleMDE({ 
    element: document.getElementById("todo-body"),
    autosave: {
        enable: false
    },
    forceSync: true,
});

document.body.addEventListener('htmx:configRequest', function(evt){
    params = evt.detail.parameters
    html = simplemde.options.previewRender(params['body'])
    evt.detail.parameters['body_encoded'] = html
    simplemde.value("")
})