import * as React from 'react'
import { DismissibleAlert } from './DismissibleAlert'

/**
 * A global alert telling all users that due to Docker for Mac, site performance
 * will be degraded.
 */
export const DockerForMacAlert: React.SFC = () => (
    <DismissibleAlert partialStorageKey="DockerForMac" className="docker-for-mac-alert">
        <span className="docker-for-mac-alert__left">
            It looks like you're using Docker for Mac. Due to known issues related to Docker for Mac's file system
            access, search performance and cloning repositories on Sourcegraph will be much slower.
        </span>
        <span className="docker-for-mac-alert__right">
            <a target="_blank" href="https://about.sourcegraph.com/docs">
                Run Sourcegraph on a different platform or deploy it to a server
            </a>{' '}
            for much faster performance.
        </span>
    </DismissibleAlert>
)
