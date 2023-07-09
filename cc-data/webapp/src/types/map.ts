export type MapData = {
    items: Country[];
    points: PointData;
};

export type Country = {
    id: string;
    iso3: string;
    coordinates?: any;
    country: string;
    value: number;
};

export type PointData = {
    defaultPoint: Point;
    points: Point[];
};

export type Point = {
    value: string;
    label: string;
};
