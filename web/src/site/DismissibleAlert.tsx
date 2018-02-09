import CloseIcon from '@sourcegraph/icons/lib/Close'
import * as React from 'react'

interface Props {
    /** used to build the key that represents the alert in local storage */
    partialStorageKey: string

    /** class name to be applied to the alert */
    className: string
}

interface State {
    dismissed: boolean
}

/**
 * A global site alert that can be dismissed. Once dismissed, it is never shown
 * again.
 */
export class DismissibleAlert extends React.PureComponent<Props, State> {
    private key: string

    constructor(props: Props) {
        super(props)
        this.key = 'DismissibleAlert/' + props.partialStorageKey + '/dismissed'

        this.state = {
            dismissed: localStorage.getItem(this.key) === 'true',
        }
    }

    public render(): JSX.Element | null {
        if (this.state.dismissed) {
            return null
        }
        return (
            <div className={`alert alert-warning site-alert dismissible-alert ${this.props.className}`}>
                {this.props.children}
                <button className="btn btn-icon" onClick={this.onDismiss}>
                    <CloseIcon />
                </button>
            </div>
        )
    }

    private onDismiss = (): void => {
        localStorage.setItem(this.key, 'true')
        this.setState({ dismissed: true })
    }
}
