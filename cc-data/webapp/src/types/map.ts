export type MapData = {
    items: Country[];
    points: PointData;
    range: number[];
    colorRange: string[];
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