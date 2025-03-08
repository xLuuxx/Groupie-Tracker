mapboxgl.accessToken = 'pk.eyJ1Ijoib2loYSIsImEiOiJjbTYweGEycHEwaTJsMmxzNTBjanU0OWVvIn0.IZDNi1lmxBbldfGE7VhSrA';
const map = new mapboxgl.Map({
    container: 'map',
    center: geojson.features[0].geometry.coordinates,
    zoom: 15.5,
    bearing: 130,
    pitch: 75,
    style: 'mapbox://styles/oiha/cm65twlsa00cz01s7bafy17wk',
});

map.addControl(new mapboxgl.FullscreenControl({
    position: 'top-right'
}));

map.on('style.load', () => {
    map.addSource('mapbox-dem', {
        'type': 'raster-dem',
        'url': 'mapbox://mapbox.mapbox-terrain-dem-v1',
        'tileSize': 512,
        'maxzoom': 14
    });
    map.setTerrain({ 'source': 'mapbox-dem', 'exaggeration': 1.5 });
});

const populateMarkersAndSelect = () => {
    const mySelect = document.getElementById('locations');

    mySelect.innerHTML = '';

    geojson.features.forEach((marker, index) => {
        const el = document.createElement('div');
        el.className = 'marker';

        new mapboxgl.Marker(el)
            .setLngLat(marker.geometry.coordinates)
            .setPopup(new mapboxgl.Popup({ offset: 25 })
                .setHTML(`<h3>${marker.properties.title}</h3><p>${marker.properties.description}</p>`))
            .addTo(map);

        const option = document.createElement('option');
        option.value = index;
        option.textContent = marker.properties.title;
        option.dataset.coordinates = JSON.stringify(marker.geometry.coordinates);
        mySelect.appendChild(option);
    });
};

populateMarkersAndSelect();

const mySelect = document.getElementById('locations');
mySelect.addEventListener('change', (event) => {
    const selectedOption = event.target.options[event.target.selectedIndex];
    const coordinates = JSON.parse(selectedOption.dataset.coordinates);

    const location = {
        center: coordinates,
        zoom: 15.5,
        bearing: 130,
        pitch: 75
    };

    map.flyTo({
        ...location,
        duration: 12000,
        essential: true
    });
});

