import { useRef, useEffect, useState } from "react";
import mapboxgl from "mapbox-gl";
import { FeatureCollection, Feature, Point } from "geojson";
import barrattLogo from "./logos/Barratt.png";
import bellwayLogo from "./logos/Bellway.png";
import berkeleyLogo from "./logos/Berkeley.png";
import persimmonLogo from "./logos/Persimmon.png";
import taylorWimpeyLogo from "./logos/TaylorWimpey.png";

// https://github.com/visgl/react-map-gl/issues/1266
// @ts-ignore
// eslint-disable-next-line import/no-webpack-loader-syntax
import MapboxWorker from "worker-loader!mapbox-gl/dist/mapbox-gl-csp-worker";

// @ts-ignore
mapboxgl.workerClass = MapboxWorker;

const builderColors: [string, string][] = [
  ["Barratt", "#a3ca36"],
  ["Bellway", "#f37021"],
  ["Berkeley", "#cc0e25"],
  ["Persimmon", "#33b28b"],
  ["Taylor Wimpey", "#735394"],
];

const builderLogos: Record<string, string> = {
  Barratt: barrattLogo,
  Bellway: bellwayLogo,
  Berkeley: berkeleyLogo,
  Persimmon: persimmonLogo,
  "Taylor Wimpey": taylorWimpeyLogo,
};

interface DevelopmentResult {
  builder: string;
  lat: number;
  lng: number;
  location: string;
  name: string;
  url: string;
}

interface DevelopmentData {
  scraped_at: string;
  results: DevelopmentResult[];
}

function App() {
  const mapContainer = useRef<HTMLDivElement | null>(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const [data, setData] = useState<DevelopmentData>();

  useEffect(() => {
    if (map.current) {
      return;
    }

    map.current = new mapboxgl.Map({
      container: mapContainer.current as HTMLElement,
      accessToken: process.env.REACT_APP_MAPBOX_ACCESS_TOKEN,
      style: "mapbox://styles/mapbox/dark-v11",
      center: [-3.4735, 54.1171],
      zoom: 4.5,
      projection: {
        name: "mercator",
      },
    });

    map.current.on("load", () => {
      if (!map.current) {
        throw new Error("could not get map ref");
      }

      (async () => {
        const res = await fetch(`${process.env.PUBLIC_URL}/developments.json`);
        const data = (await res.json()) as DevelopmentData;
        setData(data);
        const results = convertToGeoJson(data.results);

        if (!map.current) {
          throw new Error("could not get map ref");
        }

        map.current.addSource("developments", {
          type: "geojson",
          data: results,
        });

        addLayers(map.current);
        addClickListeners(map.current);
      })();
    });
  });

  return (
    <>
      <header>
        <h1>UK Housing Developments</h1>
        <a
          href="https://github.com/oliverbenns/uk-housing-developments"
          target="_blank"
          rel="noreferrer"
        >
          View on Github
        </a>
        <ul>
          {builderColors.map((val) => {
            const [name, color] = val;
            return (
              <li key={name}>
                <span className="bull" style={{ backgroundColor: color }} />
                {name}
              </li>
            );
          })}
        </ul>
        {data && <i>Data at {data.scraped_at.split("T")[0]}</i>}
      </header>
      <main>
        <div ref={mapContainer} className="map" />
      </main>
    </>
  );
}

const addLayers = (map: mapboxgl.Map) => {
  map.addLayer({
    id: "developments-point",
    type: "circle",
    source: "developments",
    paint: {
      "circle-radius": {
        base: 1.75,
        stops: [
          [2, 2],
          [8, 4],
          [12, 14],
        ],
      },
      "circle-color": [
        "match",
        ["get", "builder"],
        ...builderColors.flat(),
        "red",
      ],
    },
  });
};

const addClickListeners = (map: mapboxgl.Map) => {
  map.on("click", "developments-point", (ev) => {
    if (!ev.features) {
      return;
    }

    const feature = ev.features[0];
    const coordinates = (feature.geometry as any).coordinates.slice();
    const development = feature.properties as DevelopmentResult;
    const logoSrc = builderLogos[development.builder];
    const [, color] = builderColors.find(
      (val) => val[0] === development.builder,
    )!;

    new mapboxgl.Popup()
      .setLngLat(coordinates)
      .setHTML(
        `
          <div>
            <img src=${logoSrc} alt="${development.builder} logo"/>
            <h2>${development.name}</h2>
            <p>${development.lng.toFixed(6)} / ${development.lat.toFixed(6)}</p>

            <a href="${
              development.url
            }" class="button" target="_blank" rel="noreferrer" style="background-color: ${color}">View</a>
          </div>
      `,
      )
      .addTo(map);
  });

  // Change the cursor to a pointer when the mouse is over the places layer.
  map.on("mouseenter", "developments-point", () => {
    map.getCanvas().style.cursor = "pointer";
  });

  // Change it back to a pointer when it leaves.
  map.on("mouseleave", "developments-point", () => {
    map.getCanvas().style.cursor = "";
  });
};

const convertToGeoJson = (
  results: DevelopmentResult[],
): FeatureCollection<Point, DevelopmentResult> => {
  const features = results.map((result): Feature<Point, DevelopmentResult> => {
    return {
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [result.lng, result.lat],
      },
      properties: result,
    };
  });

  return {
    type: "FeatureCollection",
    features: features,
  };
};

export default App;
