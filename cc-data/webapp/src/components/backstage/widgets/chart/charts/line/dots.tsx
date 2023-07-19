import React, {
    FC,
    useContext,
    useEffect,
    useState,
} from 'react';
import {useIntl} from 'react-intl';
import {useRouteMatch} from 'react-router-dom';
import styled from 'styled-components';

import {Tooltip} from 'antd';

import {formatUrlAsMarkdown} from 'src/components/backstage/header/controls';
import {useToaster} from 'src/components/backstage/toast_banner';
import {FullUrlContext, IsRhsClosedContext} from 'src/components/rhs/rhs';
import {IsEcosystemRhsContext} from 'src/components/rhs/rhs_widgets';
import {
    buildQuery,
    buildTo,
    buildToForCopy,
    useUrlHash,
} from 'src/hooks';
import {copyToClipboard} from 'src/utils';
import {LineDot} from 'src/types/charts';

export const DOT_PREFIX = 'dot-';

const dotStringify = (x: string, y: string, value: string | number): [string, string, string] => {
    const cxString = `${x}`.replace('.', 'dot');
    const cyString = `${y}`.replace('.', 'dot');
    const valueString = `${value}`.replace('.', 'dot');
    return [cxString, cyString, valueString];
};

export const valueStringify = (value: string | number): string => {
    return `${value}`.replace('.', 'dot');
};

export const idStringify = (id: string): string => {
    return `${id}`.replaceAll('-', 'hpn');
};

const isDotSelected = (selectedDot: LineDot, label: string, value: number): boolean => {
    return selectedDot.label === label && valueStringify(selectedDot.value) === valueStringify(value);
};

type Props = any;

export const Dot: FC<Props> = (props) => {
    const {cx, cy, payload, value, originalColor, selectedDot, sectionId} = props;

    const isRhsClosed = useContext(IsRhsClosedContext);
    const urlHash = useUrlHash();

    const [open, setOpen] = useState(false);
    const [color, setColor] = useState(originalColor);

    useEffect(() => {
        setOpen(false);
        const timeout = setTimeout(() => {
            const isSelected = isDotSelected(selectedDot, payload.label, value);
            setOpen(isSelected);
            setColor(isSelected ? '#F4B400' : originalColor);
        }, 100);
        return () => {
            clearTimeout(timeout);
        };
    }, [isRhsClosed, selectedDot, urlHash]);

    // useEffect(() => {
    //     const isSelected = isDotSelected(selectedDot, payload.label, value);
    //     setOpen(isSelected);
    //     setColor(isSelected ? '#F4B400' : originalColor);
    // }, [selectedDot]);

    return (
        <Tooltip
            title={`${payload.label}: ${value}`}
            open={open}
        >
            <DotCircle
                id={`dot-${payload.label}-${valueStringify(value)}-${idStringify(sectionId)}`}
                cx={cx}
                cy={cy}
                r={4}
                fill={color}
            />
        </Tooltip>
    );
};

export const ClickableDot: FC<Props> = (props) => {
    const {cx, cy, payload, value, originalColor, selectedDot, parentId, sectionId} = props;

    const isEcosystemRhs = useContext(IsEcosystemRhsContext);
    const fullUrl = useContext(FullUrlContext);
    const {url} = useRouteMatch();
    const {formatMessage} = useIntl();
    const {add: addToast} = useToaster();

    const [color, setColor] = useState(originalColor);

    const ecosystemQuery = isEcosystemRhs ? '' : buildQuery(parentId, sectionId);
    const [,, valueString] = dotStringify(cx, cy, value);

    useEffect(() => {
        const isSelected = isDotSelected(selectedDot, payload.label, value);
        setColor(isSelected ? '#F4B400' : originalColor);
    }, [selectedDot]);

    const handleDotClick = (event: any) => {
        const itemId = `dot-${payload.label}-${valueString}-${idStringify(sectionId)}`;
        const name = `${payload.label}: ${value}`;
        const path = buildToForCopy(buildTo(fullUrl, itemId, ecosystemQuery, url));
        copyToClipboard(formatUrlAsMarkdown(path, name));
        addToast({content: formatMessage({defaultMessage: 'Copied!'})});
    };

    return (
        <DotCircle
            cx={cx}
            cy={cy}
            r={7}
            fill={color}
            onClick={handleDotClick}
        />
    );
};

const DotCircle = styled.circle`
    cursor: pointer;    
`;