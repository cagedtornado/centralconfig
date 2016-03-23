import React from 'react';
window.React = React; // export for http://fb.me/react-devtools

//	The API utils
import CentralConfigAPIUtils from './utils/CentralConfigAPIUtils';

//	The app component
import CentralConfigApp from './components/CentralConfigApp.react';

//	Application element
var appElement = document.getElementById("centralconfigapp");

//	Start the app
React.render(<CentralConfigApp />, appElement);