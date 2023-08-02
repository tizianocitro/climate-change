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
    isUrlHashValid,
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
import Legend, {LegendData} from './legend';

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

// Extracts the country ISO3.
// At a certain point the url hash is modified by getUrlHashForWorldMap() to provide the correct hash for useScrollIntoView().
// This is way we don't use the point in the url hash, because it has been removed by getUrlHashForWorldMap().
export const getCountryFromUrlHashWithoutPoint = (urlHash: string): string => {
    const urlHashNoPrefix = urlHash.substring(7);
    const segments = urlHashNoPrefix.split('-');
    const iso3OrPoint = segments[0];

    // If it's not 3 characters long, then it's the point, otherwise it's the ISO3
    if (iso3OrPoint.length !== 3) {
        return segments[1];
    }
    return iso3OrPoint;
};

export const getUrlHashForWorldMap = (
    urlHash: string,
    sectionId: string,
    isRhs: boolean,
): string => {
    if (!isUrlHashValid(urlHash, [sectionId], ['mapel-', 'sea-'])) {
        return urlHash;
    }
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
    if (segments.length < 2) {
        return '';
    }
    let worldMapHash = '';
    if (isRhs) {
        worldMapHash = `${prefix}${segments.slice(1).join('-')}`;
    } else { // prevent losing country ISO3 in url hash when in dashboard
        worldMapHash = `${prefix}${segments.join('-')}`;
        if (segments[0].length !== 3) {
            worldMapHash = `${prefix}${segments.slice(1).join('-')}`;
        }
    }
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

const getCountryValue = (countries: Country[], countryIso3: string): string | number => {
    let value: string | number = 'unknown';
    const countryForValue = countries.find((country) => country.iso3 === countryIso3);
    if (countryForValue) {
        value = countryForValue.value;
    }
    return value;
};

const getLegendData = (
    range: number[] | undefined,
    colorRange: string[] | undefined,
    worldEnv: WorldEnv | undefined,
    seaEnv: SeaEnv | undefined,
): LegendData | null => {
    if (range && colorRange) {
        return {
            minValue: range[0],
            maxValue: range[1],
            minColor: colorRange[0],
            maxColor: colorRange[1],
        };
    }
    if (worldEnv) {
        return {
            minValue: worldEnv.range[0],
            maxValue: worldEnv.range[1],
            minColor: worldEnv.colorRange[0],
            maxColor: worldEnv.colorRange[1],
        };
    }
    if (seaEnv) {
        return {
            minValue: seaEnv.range[0],
            maxValue: seaEnv.range[1],
            minColor: seaEnv.colorRange[0],
            maxColor: seaEnv.colorRange[1],
        };
    }
    return null;
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

    const legendData = getLegendData(range, colorRange, worldEnv, seaEnv);

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

    const getCountryFromMap = (
        countryIso3: string,
        mapEntry: string,
    ): Country | null => {
        let mapCountry = geoMap[mapEntry];
        if (mapCountry) {
            return mapCountry;
        }
        let hashedCountry = data.find((country) => country.iso3 === countryIso3);
        const geography = savedGeographies.find((geo: any) => geo.id === countryIso3);
        const hashedCountryCoordinates = calcCoordinates(
            geography?.geometry?.coordinates,
            geography?.geometry?.type,
        );
        if (!hashedCountryCoordinates) {
            return null;
        }

        let name = geography?.properties?.name ?? hashedCountry?.country ?? '';
        name = isSeaEnv ? 'Sea' : name;

        // TODO: maybe value can just be set to empty string here
        let value = isWorldEnv ? `${worldEnv?.value}` : 'unknown';
        value = isSeaEnv ? `${seaEnv?.value}` : 'unknown';
        if (!hashedCountry) {
            hashedCountry = {
                id: countryIso3,
                country: geography?.properties?.name,
                iso3: countryIso3,
                value,
            };
        }

        mapCountry = {
            ...hashedCountry,
            name,
            coordinates: hashedCountryCoordinates,
        };
        setGeoMap((prev: any) => ({...prev, [mapEntry]: mapCountry}));
        return mapCountry;
    };

    const updateSelectedCountry = (currentCountries: Country[]): void => {
        if (!isUrlHashValid(urlHash, [sectionId], ['mapel-', 'sea-'])) {
            return;
        }
        if (!savedGeographies) {
            return;
        }
        if (isSeaEnv && seaEnv) {
            const seaIso3 = 'MRT';
            const seaCountry = getCountryFromMap(seaIso3, 'sea');
            if (!seaCountry) {
                return;
            }
            setSelectedCountry({...seaCountry, value: seaEnv.value});
            return;
        }
        if (isWorldEnv && worldEnv) {
            const worldIso3 = getCountryFromUrlHashWithoutPoint(urlHash);
            const worldCountry = getCountryFromMap(worldIso3, worldIso3);
            if (!worldCountry) {
                return;
            }
            setSelectedCountry({...worldCountry, value: worldEnv.value});
            return;
        }
        const countryIso3 = getCountryFromUrlHashWithoutPoint(urlHash);
        const country = getCountryFromMap(countryIso3, countryIso3);
        if (!country) {
            return;
        }
        const value = getCountryValue(currentCountries, countryIso3);
        setSelectedCountry({...country, value});
    };

    useEffect(() => {
        setData(countries || []);
        updateSelectedCountry(countries || []);
    }, [countries]);

    useEffect(() => {
        updateSelectedCountry(data);
    }, [worldEnv, seaEnv]);

    // Adding urlHash as dependency helps a lot with performances when there is an ecosystem, but prevents hyperlinking to work in dashboard
    // To make it work also in the dashboard, we used a memoize approach with a cache for the hashed countries.
    // In this way, we do not have to calculate the coordinates all of that every single time.
    useEffect(() => {
        updateSelectedCountry(data);
    }, [urlHash]);

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
                                height='40'
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
                                {/* {selectedCountry.name} */}
                                {/* {`${selectedCountry.name}: ${isSeaEnv ? seaEnv?.value : selectedCountry.value}`} */}
                                <tspan
                                    x='0'
                                    dy='0'
                                >
                                    {selectedCountry.name}
                                </tspan>
                                <tspan
                                    x='0'
                                    dy='20'
                                >
                                    {isSeaEnv ? seaEnv?.value : selectedCountry.value}
                                </tspan>
                            </text>
                        </Annotation>
                    }
                </ZoomableGroup>
            </ComposableMap>

            {legendData && <Legend data={legendData}/>}

            <ReactTooltip id='geo-tooltip'/>
            <ReactTooltip id='sea-tooltip'/>
        </div>
    );
};

export default WorldMap;
