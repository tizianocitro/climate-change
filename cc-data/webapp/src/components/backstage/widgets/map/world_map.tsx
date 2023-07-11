import React, {
    FC,
    MouseEvent,
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
import {Country, Point} from 'src/types/map';

import Features from './features.json';

type Props = {
    data: Country[];
    range: number[];
    colorRange: string[];
    selectedPoint: Point;
    parentId: string;
    sectionId: string;
};

// Removes #_ at the start and then extracts the country ISO3
export const getCountryFromUrlHash = (urlHash: string): string => urlHash.substring(2).split('_')[1];

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

const WorldMap: FC<Props> = ({
    data,
    range,
    colorRange,
    selectedPoint,
    parentId,
    sectionId,
}) => {
    const isEcosystemRhs = useContext(IsEcosystemRhsContext);
    const isRhs = useContext(IsRhsContext);
    const fullUrl = useContext(FullUrlContext);
    const {url} = useRouteMatch();
    const ecosystemQuery = isEcosystemRhs ? '' : buildQuery(parentId, sectionId);

    const [tooltipContent, setTooltipContent] = useState('');
    const [tooltipPosition, setTooltipPosition] = useState({x: 0, y: 0});
    const [savedGeographies, setSavedGeographies] = useState<any>(null);
    const [selectedCountry, setSelectedCountry] = useState<Country | undefined>();

    const urlHash = useUrlHash();
    const {formatMessage} = useIntl();
    const {add: addToast} = useToaster();

    const colorScale = scaleLinear<string>().domain(range).range(colorRange);

    const handleMouseEnter = (event: MouseEvent, geo: any) => {
        const {name} = geo.properties;
        const countryData = data.find((country) => country.iso3 === geo.id);
        if (!countryData) {
            setTooltipContent(`${name}: unknown`);
            return;
        }
        setTooltipContent(`${name}: ${countryData.value}`);
    };

    const handleMouseLeave = () => {
        setTooltipContent('');
    };

    const handleMouseMove = ({clientX, clientY}: MouseEvent) => {
        setTooltipPosition({x: clientX, y: clientY});
    };

    const handleClick = (geo: any) => {
        const itemId = `_${selectedPoint.value}_${geo.id}`;
        const path = buildToForCopy(buildTo(fullUrl, itemId, ecosystemQuery, url));
        copyToClipboard(formatUrlAsMarkdown(path, `${geo.properties.name} [${selectedPoint.value}]`));
        addToast({content: formatMessage({defaultMessage: 'Copied!'})});
    };

    useEffect(() => {
        if (!savedGeographies || !urlHash || urlHash.length < 1) {
            return;
        }
        const countryIso3 = getCountryFromUrlHash(urlHash);
        let hashedCountry = data.find((country) => country.iso3 === countryIso3);
        const geography = savedGeographies.find((geo: any) => geo.id === countryIso3);
        const hashedCountryCoordinates = calcCoordinates(
            geography?.geometry?.coordinates,
            geography?.geometry?.type);
        if (!hashedCountryCoordinates) {
            return;
        }
        const name = geography?.properties?.name ?? hashedCountry?.country ?? '';
        if (!hashedCountry) {
            hashedCountry = {
                id: countryIso3,
                country: geography?.properties?.name,
                iso3: countryIso3,
                value: 'unknown',
            };
        }

        setSelectedCountry({
            ...hashedCountry,
            name,
            coordinates: hashedCountryCoordinates,
        });
    });

    const width = isRhs ? '100%' : '90%';
    return (
        <div
            style={{
                width,
                maxWidth: '100%',
                margin: '0 auto',
            }}
            onMouseMove={handleMouseMove}
        >
            <ComposableMap>
                <ZoomableGroup>
                    <Sphere
                        id='world-sphere'
                        stroke='#E4E5E6'
                        fill='#87CEEB'
                        strokeWidth={0.5}
                    />
                    <Graticule
                        stroke='#E4E5E6'
                        strokeWidth={0.5}
                    />
                    <Geographies geography={Features}>
                        {({geographies}) => geographies.map((geo) => {
                            if (!savedGeographies) {
                                setSavedGeographies(geographies);
                            }
                            const countryData = data.find((country) => country.iso3 === geo.id);
                            const color = countryData && typeof countryData.value === 'number' ?
                                colorScale(countryData.value) :
                                '#dddddd';

                            return (
                                <Geography
                                    key={geo.rsmKey}
                                    geography={geo}
                                    fill={color}
                                    stroke='#000000'
                                    style={{
                                        hover: {
                                            cursor: 'pointer',
                                        },
                                    }}
                                    onClick={() => handleClick(geo)}
                                    onMouseEnter={(event) => handleMouseEnter(event, geo)}
                                    onMouseLeave={handleMouseLeave}
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

            {tooltipContent &&
                <div
                    style={{
                        position: 'absolute',
                        top: tooltipPosition.y - 40,
                        left: tooltipPosition.x - 5,
                        background: '#FFFFFF',
                        padding: '5px 10px',
                        borderRadius: 5,
                        pointerEvents: 'none',
                        zIndex: 9999,
                    }}
                >
                    {tooltipContent}
                </div>
            }
        </div>
    );
};

export default WorldMap;
