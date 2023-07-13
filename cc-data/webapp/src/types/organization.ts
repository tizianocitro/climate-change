import {ChartType, WidgetType} from 'src/components/backstage/widgets/widget_types';

export interface PlatformConfig {
    organizations: Organization[];
}

export interface Organization {
    isEcosystem?: boolean,
    description?: string,
    id: string,
    name: string,
    sections: Section[],
    widgets: Widget[],
}

export interface Section {
    id: string,
    internal: boolean;
    name: string,
    url: string,
    sections: Section[],
    widgets: Widget[],
}

export interface Widget {
    name?: string,
    type: WidgetType,
    url?: string,
    chartType?: ChartType;
}

export interface SectionInfo {
    id: string;
    name: string;
    [key: string]: any;
}

export interface SectionInfoParams {
    name: string;
    [key: string]: any;
}

export interface FetchOrganizationsParams {
    direction?: string;
    page: number;
    per_page: number;
    organization_id?: string;
    search_term?: string;
    sort?: string;
    team_id?: string;
}

export type Object = any;