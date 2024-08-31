(function () {
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
      const { name, address, city, state, country, website, phoneNumber } = loc;
      const markerOptions = loc.isCo404Loc ? { icon: redIcon } : {};
      const args = [loc.latitude, loc.longitude];
  
      const formattedWebsite = website ? website.replace(/^(https?:\/\/|\/\/)/, '') : null;
  
      const popupContent = `
        <div class="popup-template">
          <h3 class="text-lg font-bold">${name}</h3>
          <div class="text-md">
            <div>${address}</div>
            <div>${city}${!!state ? ' ' : ''}${state}, ${country}</div>
            <div class="${!website ? 'hidden' : ''}">
              <a href="https://${formattedWebsite}" target="_blank" rel="noopener noreferrer">${website}</a>
            </div>
            <div class="${!phoneNumber ? 'hidden' : ''}">
              <a href="tel:${phoneNumber}" target="_blank" rel="noopener noreferrer">${phoneNumber}</a>
            </div>
          </div>
        </div>
      `;
    
      const marker = L.marker(args, markerOptions).bindPopup(popupContent);
      markers.addLayer(marker);
    });
  
    map.addLayer(markers);
    
    const group = new L.featureGroup(markers.getLayers());
    map.fitBounds(group.getBounds());
  });
  
})() 
