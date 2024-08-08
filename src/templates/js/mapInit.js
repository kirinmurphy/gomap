import htmx from 'htmx.org';
import L from 'leaflet';

const map = L.map('map').setView([0, 0], 2);


console.log('fffffffff'); 
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);

const redIcon = new L.Icon({
  iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
  className: 'red-marker' // Add a custom class to apply CSS
});
    
htmx.on('htmx:afterSettle', function(evt) {
  const locations = JSON.parse(evt.detail.xhr.response);
  locations.forEach(function(loc) {
    const markerOptions = loc.isCo404Loc ? { icon: redIcon } : {};  
    L.marker([loc.latitude, loc.longitude], markerOptions).addTo(map)
      .bindPopup(loc.name);
  });
});
