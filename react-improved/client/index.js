import React from 'react'
import ReactDOM from 'react-dom'
import App from './components/App.js'
import { NotieProvider } from 'react-notie'

const EntryPoint = (
    <NotieProvider position='bottom'>
        <App/>
    </NotieProvider>
)

ReactDOM.render(EntryPoint, document.getElementById('react-app'))
