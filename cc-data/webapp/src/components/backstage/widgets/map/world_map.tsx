import React, {
    FC,
    useContext,
    useEffect,
    useState,
} from 'react';
import {useIntl} from 'react-intl';
import {
    Annotation,
    ComposableMap,
    Geographies,
    Geography,
    Graticule,
    Sphere,
    ZoomableGroup,
} from 'react-simple-maps';
import {useRouteMatch} from 'react-router-dom';
import {scaleLinear} from 'd3-scale';
import {Tooltip as ReactTooltip} from 'react-tooltip';

import {
    buildQuery,
    buildTo,
    buildToForCopy,
    useUrlHash,
} from 'src/hooks';
import {copyToClipboard} from 'src/utils';
import {formatUrlAsMarkdown} from 'src/components/backstage/header/controls';
import {useToaster} from 'src/components/backstage/toast_banner';
import {IsEcosystemRhsContext} from 'src/components/rhs/rhs_widgets';
import {IsRhsContext} from 'src/components/backstage/sections_widgets/sections_widgets_container';
import {FullUrlContext} from 'src/components/rhs/rhs';
import {
    Country,
    Point,
    SeaEnv,
    WorldEnv,
    isSeaEnvDefined,
    isWorldEnvDefined,
} from 'src/types/map';

import Features from './features.json';

type Props = {
    countries?: Country[];
    range?: number[];
    colorRange?: string[];
    worldEnv?: WorldEnv;
    seaEnv?: SeaEnv;
    selectedPoint: Point;
    parentId: string;
    sectionId: string;
};

const defaultCountriesColor = '#8B4513';
const defaultSeaColor = '#87CEEB';

// Removes #_ at the start and then extracts the country ISO3
export const getCountryFromUrlHash = (urlHash: string): string => urlHash.substring(7).split('-')[1];

export const getUrlHashForWorldMap = (urlHash: string): string => {
    let prefix = '';
    let segments: string[] = [];
    if (urlHash.includes('mapel-')) {
        prefix = urlHash.substring(0, 7);
        segments = urlHash.substring(7).split('-');
    }
    if (urlHash.includes('sea-')) {
        prefix = urlHash.substring(0, 5);
        segments = urlHash.substring(5).split('-');
    }
    const worldMapHash = `${prefix}${segments.slice(1).join('-')}`;
    return prefix !== '' && segments.length > 0 ? worldMapHash : urlHash;
};

const calcCoordinates = (coordinates: any, type: string): any | null => {
    if (!coordinates) {
        return null;
    }

    let lat = 0;
    let long = 0;
    let count = 0;
    coordinates.forEach((coordinate: any) => {
        coordinate.forEach((polygonOrMultiPolygon: any) => {
            if (type === 'Polygon') {
                lat += parseFloat(polygonOrMultiPolygon[0]);
                long += parseFloat(polygonOrMultiPolygon[1]);
                count++;
            } else {
                polygonOrMultiPolygon.forEach((multiPolygon: any) => {
                    lat += parseFloat(multiPolygon[0]);
                    long += parseFloat(multiPolygon[1]);
                    count++;
                });
            }
        });
    });
    return [lat / count, long / count];
};

const getContriesColor = (seaEnv: SeaEnv | undefined): string => {
    return seaEnv?.countriesColor ? seaEnv?.countriesColor : defaultCountriesColor;
};

const WorldMap: FC<Props> = ({
    countries,
    range,
    colorRange,
    worldEnv,
    seaEnv,
    selectedPoint,
    parentId,
    sectionId,
}) => {
    const isEcosystemRhs = useContext(IsEcosystemRhsContext);
    const isRhs = useContext(IsRhsContext);
    const fullUrl = useContext(FullUrlContext);
    const {url} = useRouteMatch();
    const ecosystemQuery = isEcosystemRhs ? '' : buildQuery(parentId, sectionId);

    const [savedGeographies, setSavedGeographies] = useState<any>(null);
    const [selectedCountry, setSelectedCountry] = useState<Country | undefined>();
    const [geoMap, setGeoMap] = useState<any>({});
    const [data, setData] = useState<Country[]>(countries || []);
    useEffect(() => {
        setData(countries || []);
    }, [countries]);

    const urlHash = useUrlHash();
    const {formatMessage} = useIntl();
    const {add: addToast} = useToaster();

    const countriesColorScale = scaleLinear<string>().domain(range ?? []).range(colorRange ?? []);
    const seaColorScale = scaleLinear<string>().
        domain(seaEnv?.range ?? []).
        range(seaEnv?.colorRange ?? []);
    const worldColorScale = scaleLinear<string>().
        domain(worldEnv?.range ?? []).
        range(worldEnv?.colorRange ?? []);

    const countriesColor = getContriesColor(seaEnv);
    const seaColor = seaEnv ? seaColorScale(seaEnv.value) : defaultSeaColor;

    const isWorldEnv = isWorldEnvDefined(worldEnv);
    const isSeaEnv = isSeaEnvDefined(seaEnv);
    const isCountriesEnv = !isWorldEnv && !isSeaEnv;

    const getSeaTooltipContent = (): string => {
        if (!seaEnv) {
            return '';
        }
        const label = seaEnv.label || 'Sea';
        return `${label}: ${seaEnv.value}`;
    };

    const getEarthTooltipContent = (geo: any): string => {
        if (seaEnv && (seaEnv.noCountriesValue || false)) {
            return '';
        }
        const {name} = geo.properties;
        const countryData = data.find((country) => country.iso3 === geo.id);
        if (countryData) {
            return `${name}: ${countryData.value}`;
        }
        return isWorldEnv ? `${name}: ${worldEnv?.value}` : `${name}: unknown`;
    };

    const handleSeaClick = () => {
        if (!seaEnv) {
            return;
        }
        const label = seaEnv?.label ?? 'sea';
        const itemId = `sea-${selectedPoint.value}-${sectionId}`;
        const path = buildToForCopy(buildTo(fullUrl, itemId, ecosystemQuery, url));
        copyToClipboard(formatUrlAsMarkdown(path, `${label} [${selectedPoint.value}]`));
        addToast({content: formatMessage({defaultMessage: 'Copied!'})});
    };

    const handleEarthClick = (geo: any) => {
        if (isSeaEnv) {
            return;
        }
        const itemId = `mapel-${selectedPoint.value}-${geo.id}-${sectionId}`;
        const path = buildToForCopy(buildTo(fullUrl, itemId, ecosystemQuery, url));
        copyToClipboard(formatUrlAsMarkdown(path, `${geo.properties.name} [${selectedPoint.value}]`));
        addToast({content: formatMessage({defaultMessage: 'Copied!'})});
    };

    useEffect(() => {
        if (!savedGeographies || !urlHash || urlHash.length < 1 || !urlHash.includes(sectionId)) {
            return;
        }
        if (!isSeaEnv) {
            return;
        }
        const countryIso3 = 'MRT';
        let mapCountry = geoMap.sea;
        if (!mapCountry) {
            let hashedCountry = data.find((country) => country.iso3 === countryIso3);
            const geography = savedGeographies.find((geo: any) => geo.id === countryIso3);
            const hashedCountryCoordinates = calcCoordinates(
                geography?.geometry?.coordinates,
                geography?.geometry?.type,
            );
            if (!hashedCountryCoordinates) {
                return;
            }
            if (!hashedCountry) {
                hashedCountry = {
                    id: countryIso3,
                    country: geography?.properties?.name,
                    iso3: countryIso3,
                    value: isWorldEnv ? `${worldEnv?.value}` : 'unknown',
                };
            }
            const [lat, long] = hashedCountryCoordinates;
            mapCountry = {
                ...hashedCountry,
                name: 'Sea',
                coordinates: [lat, long - 30],
            };
            setGeoMap((prev: any) => ({...prev, sea: mapCountry}));
        }
        setSelectedCountry(mapCountry);

        // Adding urlHash as dependency helps a lot with performances when there is an ecosystem, but prevents hyperlinking to work in dashboard
        // To make it work also in the dashboard, we used a memoize approach with a cache for the hashed countries.
        // In this way, we do not have to calculate the coordinates all of that every single time.
    });

    useEffect(() => {
        if (!savedGeographies || !urlHash || urlHash.length < 1 || !urlHash.includes(sectionId)) {
            return;
        }
        if (isSeaEnv) {
            return;
        }
        const countryIso3 = getCountryFromUrlHash(urlHash);
        let mapCountry = geoMap[countryIso3];
        if (!mapCountry) {
            let hashedCountry = data.find((country) => country.iso3 === countryIso3);
            const geography = savedGeographies.find((geo: any) => geo.id === countryIso3);
            const hashedCountryCoordinates = calcCoordinates(
                geography?.geometry?.coordinates,
                geography?.geometry?.type,
            );
            if (!hashedCountryCoordinates) {
                return;
            }
            const name = geography?.properties?.name ?? hashedCountry?.country ?? '';
            if (!hashedCountry) {
                hashedCountry = {
                    id: countryIso3,
                    country: geography?.properties?.name,
                    iso3: countryIso3,
                    value: isWorldEnv ? `${worldEnv?.value}` : 'unknown',
                };
            }
            mapCountry = {
                ...hashedCountry,
                name,
                coordinates: hashedCountryCoordinates,
            };

            // setGeoMap((prev: any) => ({...prev, [countryIso3]: mapCountry}));
            geoMap[countryIso3] = mapCountry;
        }
        setSelectedCountry(mapCountry);

        // Adding urlHash as dependency helps a lot with performances when there is an ecosystem, but prevents hyperlinking to work in dashboard
        // To make it work also in the dashboard, we used a memoize approach with a cache for the hashed countries.
        // In this way, we do not have to calculate the coordinates all of that every single time.
    });

    const width = isRhs ? '100%' : '90%';
    return (
        <div
            style={{
                width,
                maxWidth: '100%',
                margin: '0 auto',
            }}
        >
            <ComposableMap>
                <ZoomableGroup>
                    <Sphere
                        id='world-sphere'
                        stroke='#E4E5E6'
                        fill={seaColor}
                        strokeWidth={0.5}
                    />
                    <Graticule
                        id={`sea-${sectionId}`}
                        stroke='#E4E5E6'
                        strokeWidth={0.5}
                        data-tooltip-id='sea-tooltip'
                        data-tooltip-content={getSeaTooltipContent()}
                        onClick={() => handleSeaClick()}
                        style={{cursor: isSeaEnv ? 'pointer' : 'auto'}}
                    />
                    <Geographies geography={Features}>
                        {({geographies}) => geographies.map((geo) => {
                            if (!savedGeographies) {
                                setSavedGeographies(geographies);
                            }
                            const countryData = data.find((country) => country.iso3 === geo.id);
                            let color = isCountriesEnv && countryData && typeof countryData.value === 'number' ?
                                countriesColorScale(countryData.value) :
                                countriesColor;
                            color = isWorldEnv && worldEnv ? worldColorScale(worldEnv.value) : color;

                            return (
                                <Geography
                                    id={`mapel-${geo.id}-${sectionId}`}
                                    key={geo.rsmKey}
                                    geography={geo}
                                    fill={color}
                                    stroke='#000000'
                                    style={{
                                        hover: {
                                            cursor: isSeaEnv ? 'auto' : 'pointer',
                                        },
                                    }}
                                    data-tooltip-id='geo-tooltip'
                                    data-tooltip-content={getEarthTooltipContent(geo)}
                                    onClick={() => handleEarthClick(geo)}
                                />
                            );
                        })}
                    </Geographies>
                    {selectedCountry &&
                        <Annotation
                            subject={selectedCountry.coordinates || [0, 0]}
                            dx={-70}
                            dy={-30}
                            connectorProps={{
                                stroke: '#000000',
                                strokeWidth: 1,
                                strokeLinecap: 'round',
                            }}
                        >
                            <rect
                                x={!selectedCountry.name || selectedCountry.name.length <= 14 ? '-50' : '-100'}
                                y='-11'
                                width={!selectedCountry.name || selectedCountry.name.length <= 14 ? '100' : '200'}
                                height='20'
                                rx='5'
                                ry='5'
                                fill='#FFFFFF'
                                stroke='#000000'
                                strokeWidth='0.5'
                            />
                            <text
                                textAnchor='middle'
                                alignmentBaseline='middle'
                                style={{fill: '#000000'}}
                            >
                                {selectedCountry.name}
                            </text>
                        </Annotation>
                    }
                </ZoomableGroup>
            </ComposableMap>

            <ReactTooltip id='geo-tooltip'/>
            <ReactTooltip id='sea-tooltip'/>
        </div>
    );
};

export default WorldMap;
