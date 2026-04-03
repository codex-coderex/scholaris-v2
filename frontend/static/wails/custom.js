// empty placeholder for custom.js code


// this file is here because the wails runtime expects a custom.js file to be present at /wails/custom.js. 
// you can find the function call at "frontend/node_modules/@wailsio/runtime/dist/index.js" , 
// which is the main entry point for the wails runtime.

/** exact function call 
Load custom.js if available (used by server mode for WebSocket events, etc.)
loadOptionalScript('/wails/custom.js');
**/