const map = L.map('map', {
  maxZoom: 18
}).setView([0, 0], 2);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {}).addTo(map);

const redIcon = new L.Icon({
  iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
  className: 'red-marker' // Add a custom class to apply CSS
});

const markers = L.markerClusterGroup();

htmx.on('htmx:afterSettle', function(evt) {
  const locations = JSON.parse(evt.detail.xhr.response);
  
  locations.forEach(function(loc) {
    const markerOptions = loc.isCo404Loc ? { icon: redIcon } : {};
    const args = [loc.latitude, loc.longitude];
    const marker = L.marker(args, markerOptions).bindPopup(loc.name);
    markers.addLayer(marker);
  });

  map.addLayer(markers);
  
  const group = new L.featureGroup(markers.getLayers());
  map.fitBounds(group.getBounds());
});
