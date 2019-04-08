import {
    FileSpec,
    PositionSpec,
    RepoSpec,
    RevSpec,
    toPrettyBlobURL,
    ViewStateSpec,
} from '../../../../../shared/src/util/url'
import { sourcegraphUrl } from './context'

/**
 * Returns an absolute URL to the blob (file) on the Sourcegraph instance.
 */
export function toAbsoluteBlobURL(
    ctx: RepoSpec & RevSpec & FileSpec & Partial<PositionSpec> & Partial<ViewStateSpec>
): string {
    return `${sourcegraphUrl.replace(/\/$/, '')}/${toPrettyBlobURL(ctx)}`
}
