export type MapData = {
    items?: Country[];
    points: PointData;
    range?: number[];
    colorRange?: string[];
    worldEnv?: WorldEnv;
    seaEnv?: SeaEnv;
};

export type Country = {
    id: string;
    iso3: string;
    country: string;
    value: number | string;
    name?: string;
    coordinates?: any;
};

export type PointData = {
    defaultPoint: Point;
    points: Point[];
};

export type Point = {
    value: string;
    label: string;
};

export type WorldEnv = {
    value: number;
    range: number[];
    colorRange: string[];
};

export type SeaEnv = {
    label?: string;
    value: number;
    countriesColor?: string;
    noCountriesValue?: boolean;
    range: number[];
    colorRange: string[];
};

export const defaultMapData: MapData = {
    items: [],
    points: {
        defaultPoint: {
            label: '2022',
            value: '2022',
        },
        points: [],
    },
    range: [0, 0],
    colorRange: ['', ''],
};

export const isWorldEnvDefined = (worldEnv: WorldEnv | undefined): boolean => {
    if (!worldEnv) {
        return false;
    }
    const {colorRange, range} = worldEnv;
    if (!colorRange || colorRange.length < 1) {
        return false;
    }
    return range && range.length > 0;
};

export const isSeaEnvDefined = (seaEnv: SeaEnv | undefined): boolean => {
    return seaEnv ?
        typeof seaEnv.countriesColor !== 'undefined' :
        false;
};